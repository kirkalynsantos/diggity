package output

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	log "github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/output/cyclonedx"
	"github.com/carbonetes/diggity/internal/output/github"
	"github.com/carbonetes/diggity/internal/output/save"
	"github.com/carbonetes/diggity/internal/output/spdx"
	"github.com/carbonetes/diggity/internal/output/tabular"
	"github.com/carbonetes/diggity/internal/parser/bom"
	"github.com/carbonetes/diggity/internal/parser/distro"
	"github.com/carbonetes/diggity/internal/parser/docker"
	"github.com/carbonetes/diggity/internal/secret"

	"golang.org/x/exp/maps"
)

type result map[string]*model.Package

// Result interface
var Result result = make(map[string]*model.Package, 0)

// PrintResults prints the result based on the arguments
func PrintResults() {
	finalizeResults()
	outputTypes := strings.ToLower(bom.Arguments.Output.ToOutput())

	// Table Output(Default)
	selectOutputType(outputTypes)

	if len(bom.Errors) > 0 {
		for _, err := range bom.Errors {
			log.GetLogger().Printf("[warning]: %+v\n", *err)
		}
	}
}

// Select Output Type based on the User Input with aliases considered
func selectOutputType(outputTypes string) {
	for _, output := range strings.Split(outputTypes, ",") {
		switch output {
		case model.Table:
			tabular.PrintTable()
		case model.JSON.ToOutput():
			if len(*bom.Arguments.OutputFile) > 0 {
				save.ResultToFile(GetResults())
			} else {
				fmt.Printf("%+v\n", GetResults())
			}
		case model.CycloneDXXML, "cyclonedxxml", "cyclonedx", "cyclone":
			cyclonedx.PrintCycloneDXXML()
		case model.CycloneDXJSON, "cyclonedxjson":
			cyclonedx.PrintCycloneDXJSON()
		case model.SPDXJSON, "spdxjson":
			spdx.PrintSpdxJSON()
		case model.SPDXTagValue, "spdxtagvalue", "spdx", "spdxtv":
			spdx.PrintSpdxTagValue()
		case model.GithubJSON, "githubjson", "github":
			github.PrintGithubJSON()
		}
	}
}

// Remove Duplicates and Sort Results
func finalizeResults() {
	for _, _package := range bom.Packages {
		if _, exists := Result[_package.Name+":"+_package.Version+":"+_package.Type]; !exists {
			Result[_package.Name+":"+_package.Version+":"+_package.Type] = _package
		} else {
			idx := 0
			if len(_package.Locations) > 0 {
				idx = len(_package.Locations) - 1
				for _, l := range _package.Locations {
					if l != _package.Locations[idx] {
						_package.Locations = append(_package.Locations, model.Location{
							Path:      _package.Path,
							LayerHash: "sha256:" + _package.Locations[idx].LayerHash,
						})
						Result[_package.Name+":"+_package.Version+":"+_package.Type] = _package
					}
				}
			}
		}
	}
	sortResults()
}

// Sort Results
func sortResults() {
	bom.Packages = maps.Values(Result)
	sort.Slice(bom.Packages, func(i, j int) bool {
		if bom.Packages[i].Name == bom.Packages[j].Name {
			return bom.Packages[i].Version < bom.Packages[j].Version
		}
		return bom.Packages[i].Name < bom.Packages[j].Name
	})
}

// GetResults - For event bus handler
func GetResults() string {
	_packages := maps.Values(Result)

	sort.Slice(_packages, func(i, j int) bool {
		return _packages[i].Name < _packages[j].Name
	})

	output := Output{
		Distro:   distro.Distro(),
		Packages: bom.Packages,
	}

	if !*bom.Arguments.DisableSecretSearch {
		output.Secret = secret.SecretResults
	}

	output.ImageInfo = docker.ImageInfo

	result, _ := json.MarshalIndent(output, "", " ")
	return string(result)
}
