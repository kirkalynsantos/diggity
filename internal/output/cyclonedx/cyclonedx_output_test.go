package cyclonedx

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/model/output"
	"github.com/carbonetes/diggity/internal/parser/bom"
)

type (
	AddDistroComponenetResult struct {
		_distro  *model.Distro
		expected output.Component
	}

	ConvertToComponentResult struct {
		_package *model.Package
		expected output.Component
	}

	InitPropertiesResult struct {
		_package *model.Package
		expected []output.Property
	}

	AddIDResult struct {
		_package *model.Package
		expected string
	}

	ConvertLicenseResult struct {
		_package *model.Package
		expected []output.License
	}
)

var (
	cdxPackage1 = model.Package{
		Name:    "zlib",
		Type:    "apk",
		Version: "1.2.12-r3",
		Path:    filepath.Join("lib", "apk", "db", "installed"),
		Locations: []model.Location{
			{
				Path:      filepath.Join("lib", "apk", "db", "installed"),
				LayerHash: "9b7240956cfbfefddcd91a2195bfb2ed2cd17bdff81f21111849d643dfaf8131",
			},
		},
		Description: "compression/decompression Library",
		Licenses: []string{
			"Zlib",
		},
		CPEs: []string{
			"cpe:2.3:a:zlib:zlib:1.2.12-r3:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:alpine/zlib@1.2.12-r3?arch=x86_64\u0026upstream=zlib\u0026distro="),
	}

	cdxPackage2 = model.Package{
		Name:    "libapt-pkg6.0",
		Type:    "deb",
		Version: "2.2.4",
		Path:    "",
		Locations: []model.Location{
			{
				Path:      filepath.Join("var", "lib", "dpkg", "status"),
				LayerHash: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0",
			},
			{
				Path:      filepath.Join("var", "lib", "dpkg", "status"),
				LayerHash: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0",
			},
		},
		Description: "package management runtime library This library provides the common functionality for searching and managing packages as well as information about packages.",
		Licenses: []string{
			"GPLv2+",
		},
		CPEs: []string{
			"cpe:2.3:a:libapt-pkg6.0:libapt-pkg6.0:2.2.4:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:deb/libapt-pkg6.0@2.2.4arch=s390x"),
	}
	cdxPackage3 = model.Package{
		Name:    "hardlink",
		Type:    "rpm",
		Version: "1.0",
		Path:    filepath.Join("var", "lib", "rpm", "Packages"),
		Locations: []model.Location{
			{
				Path:      filepath.Join("var", "lib", "rpm", "Packages"),
				LayerHash: "d1fd2cca7a7751ca9786b088cf639e65088fa0bda34492bb5ba292c32195461a",
			},
		},
		Description: "Create a tree of hardlinks",
		Licenses: []string{
			"GPL+",
		},
		CPEs: []string{
			"cpe:2.3:a:redhat:hardlink:1.0-19.el7:*:*:*:*:*:*:*",
			"cpe:2.3:a:hardlink:hardlink:1.0-19.el7:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:rpm/hardlink@1.0arch=x86_64"),
	}
	cdxPackage4 = model.Package{
		Name:    "phpdocumentor/reflection",
		Type:    "php",
		Version: "5.2.0",
		Path:    "phpdocumentor/reflection",
		Locations: []model.Location{
			{
				Path:      filepath.Join("opt", "phpdoc", "composer.lock"),
				LayerHash: "12a3251e94a5184b3c5f4efbc0c8df91cf8479af3745941c9d9102298d258b83",
			},
		},
		Description: "Reflection library to do Static Analysis for PHP Projects",
		Licenses: []string{
			"MIT",
		},
		CPEs: []string{
			"cpe:2.3:a:phpdocumentor:reflection:5.2.0:*:*:*:*:*:*:*",
			"cpe:2.3:a:reflection:reflection:5.2.0:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:composer/phpdocumentor/reflection@5.2.0"),
	}

	cdxDistro1 = model.Distro{
		PrettyName:   "Alpine Linux v3.16",
		Name:         "Alpine Linux",
		ID:           "alpine",
		VersionID:    "3.16.2",
		HomeURL:      "https://alpinelinux.org/",
		BugReportURL: "https://gitlab.alpinelinux.org/alpine/aports/-/issues",
	}

	cdxDistro2 = model.Distro{
		PrettyName:   "Debian GNU/Linux 11 (bullseye)",
		Name:         "Debian GNU/Linux",
		ID:           "debian",
		Version:      "11 (bullseye)",
		VersionID:    "11",
		HomeURL:      "https://www.debian.org/",
		SupportURL:   "https://www.debian.org/support",
		BugReportURL: "https://bugs.debian.org/",
	}

	cdxDistro3 = model.Distro{
		PrettyName: "CentOS Linux 8",
		Name:       "CentOS Linux",
		ID:         "centos",
		IDLike: []string{
			"rhel",
			"fedora",
		},
		Version:      "8",
		VersionID:    "8",
		HomeURL:      "https://centos.org/",
		BugReportURL: "https://bugs.centos.org/",
	}
)

func TestConvertPackage(t *testing.T) {
	bom.Packages = []*model.Package{&cdxPackage1, &cdxPackage2, &cdxPackage3, &cdxPackage4}
	expected := output.CycloneFormat{
		XMLNS: XMLN,
		Metadata: &output.Metadata{
			Tools: &[]output.Tool{
				{
					Vendor: vendor,
					Name:   name,
				},
			},
		},
		Components: &[]output.Component{
			{
				BOMRef:     "pkg:rpm/hardlink@1.0arch=x86_64?package-id=",
				Type:       library,
				Name:       "hardlink",
				Version:    "1.0",
				PackageURL: string(model.PURL("pkg:rpm/hardlink@1.0arch=x86_64")),
				Licenses: &[]output.License{
					{ID: "GPL+"},
				},
				Properties: &[]output.Property{
					{Name: "diggity:package:type", Value: "rpm"},
					{Name: "diggity:cpe23", Value: "cpe:2.3:a:redhat:hardlink:1.0-19.el7:*:*:*:*:*:*:*"},
					{Name: "diggity:cpe23", Value: "cpe:2.3:a:hardlink:hardlink:1.0-19.el7:*:*:*:*:*:*:*"},
					{Name: "diggity:location:0:layerHash", Value: "d1fd2cca7a7751ca9786b088cf639e65088fa0bda34492bb5ba292c32195461a"},
					{Name: "diggity:location:0:path", Value: `var` + string(os.PathSeparator) + `lib` + string(os.PathSeparator) + `rpm` + string(os.PathSeparator) + `Packages`},
				},
			},
			{
				BOMRef:     "pkg:deb/libapt-pkg6.0@2.2.4arch=s390x?package-id=",
				Type:       library,
				Name:       "libapt-pkg6.0",
				Version:    "2.2.4",
				PackageURL: string(model.PURL("pkg:deb/libapt-pkg6.0@2.2.4arch=s390x")),
				Licenses: &[]output.License{
					{ID: "GPLv2+"},
				},
				Properties: &[]output.Property{
					{Name: "diggity:package:type", Value: "deb"},
					{Name: "diggity:cpe23", Value: "cpe:2.3:a:libapt-pkg6.0:libapt-pkg6.0:2.2.4:*:*:*:*:*:*:*"},
					{Name: "diggity:location:0:layerHash", Value: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0"},
					{Name: "diggity:location:0:path", Value: `var` + string(os.PathSeparator) + `lib` + string(os.PathSeparator) + `dpkg` + string(os.PathSeparator) + `status`},
					{Name: "diggity:location:1:layerHash", Value: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0"},
					{Name: "diggity:location:1:path", Value: `var` + string(os.PathSeparator) + `lib` + string(os.PathSeparator) + `dpkg` + string(os.PathSeparator) + `status`},
				},
			},
			{
				BOMRef:     "pkg:composer/phpdocumentor/reflection@5.2.0?package-id=",
				Type:       library,
				Name:       "phpdocumentor/reflection",
				Version:    "5.2.0",
				PackageURL: string(model.PURL("pkg:composer/phpdocumentor/reflection@5.2.0")),
				Licenses: &[]output.License{
					{ID: "MIT"},
				},
				Properties: &[]output.Property{
					{Name: "diggity:package:type", Value: "php"},
					{Name: "diggity:cpe23", Value: "cpe:2.3:a:phpdocumentor:reflection:5.2.0:*:*:*:*:*:*:*"},
					{Name: "diggity:cpe23", Value: "cpe:2.3:a:reflection:reflection:5.2.0:*:*:*:*:*:*:*"},
					{Name: "diggity:location:0:layerHash", Value: "12a3251e94a5184b3c5f4efbc0c8df91cf8479af3745941c9d9102298d258b83"},
					{Name: "diggity:location:0:path", Value: `opt` + string(os.PathSeparator) + `phpdoc` + string(os.PathSeparator) + `composer.lock`},
				},
			},
			{
				BOMRef:     "pkg:alpine/zlib@1.2.12-r3?arch=x86_64&upstream=zlib&distro=?package-id=",
				Type:       library,
				Name:       "zlib",
				Version:    "1.2.12-r3",
				PackageURL: string(model.PURL("pkg:alpine/zlib@1.2.12-r3?arch=x86_64\u0026upstream=zlib\u0026distro=")),
				Licenses: &[]output.License{
					{ID: "Zlib"},
				},
				Properties: &[]output.Property{
					{Name: "diggity:package:type", Value: "apk"},
					{Name: "diggity:cpe23", Value: "cpe:2.3:a:zlib:zlib:1.2.12-r3:*:*:*:*:*:*:*"},
					{Name: "diggity:location:0:layerHash", Value: "9b7240956cfbfefddcd91a2195bfb2ed2cd17bdff81f21111849d643dfaf8131"},
					{Name: "diggity:location:0:path", Value: `lib` + string(os.PathSeparator) + `apk` + string(os.PathSeparator) + `db` + string(os.PathSeparator) + `installed`},
				},
			},
		},
	}

	_output := convertPackage()

	if _output.XMLNS != expected.XMLNS {
		t.Errorf("Test Failed: Expected output of %v, received: %v ", expected.XMLNS, _output.XMLNS)
	}

	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	if !r.MatchString(_output.SerialNumber) {
		t.Errorf("Test Failed: Output %v must be a valid UUID.", _output.SerialNumber)
	}

	if _output := convertPackage(); _output.Metadata.Timestamp != time.Now().Format(time.RFC3339) {
		t.Errorf("Test Failed: Must produce current timestamp of %v, received: %v.",
			time.Now().Format(time.RFC3339), getFromSource().Timestamp)
	}

	outputTools := *_output.Metadata.Tools
	if outputTools[0].Vendor != vendor ||
		outputTools[0].Name != name {
		t.Errorf("Test Failed: Expected output of %v, received: %v ", output.Tool{
			Vendor: vendor,
			Name:   name,
		}, outputTools)
	}

	if (len(*_output.Components) - 1) != len(*expected.Components) { // No distro parsed for unit test
		t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(*expected.Components), len(*_output.Components)-1)
	}
}

func TestAddDistroComponent(t *testing.T) {
	tests := []AddDistroComponenetResult{
		{&cdxDistro1, output.Component{
			Type:        operatingSystem,
			Name:        "alpine",
			Description: "Alpine Linux v3.16",
			ExternalReferences: &[]output.ExternalReference{
				{URL: "https://gitlab.alpinelinux.org/alpine/aports/-/issues", Type: "issue-tracker"},
				{URL: "https://alpinelinux.org/", Type: "website"},
			},
			Properties: &[]output.Property{
				{Name: "diggity:distro:id", Value: "alpine"},
				{Name: "diggity:distro:prettyName", Value: "Alpine Linux v3.16"},
				{Name: "diggity:distro:distributionCodename", Value: ""},
				{Name: "diggity:distro:versionID", Value: "3.16.2"},
			},
		}},
		{&cdxDistro2, output.Component{
			Type:        operatingSystem,
			Name:        "debian",
			Description: "Debian GNU/Linux 11 (bullseye)",
			ExternalReferences: &[]output.ExternalReference{
				{URL: "https://bugs.debian.org/", Type: "issue-tracker"},
				{URL: "https://www.debian.org/", Type: "website"},
				{URL: "https://www.debian.org/support", Type: "other", Comment: "support"},
			},
			Properties: &[]output.Property{
				{Name: "diggity:distro:id", Value: "debian"},
				{Name: "diggity:distro:prettyName", Value: "Debian GNU/Linux 11 (bullseye)"},
				{Name: "diggity:distro:distributionCodename", Value: ""},
				{Name: "diggity:distro:versionID", Value: "11"},
			},
		}},
		{&cdxDistro3, output.Component{
			Type:        operatingSystem,
			Name:        "centos",
			Description: "CentOS Linux 8",
			ExternalReferences: &[]output.ExternalReference{
				{URL: "https://bugs.centos.org/", Type: "issue-tracker"},
				{URL: "https://centos.org/", Type: "website"},
			},
			Properties: &[]output.Property{
				{Name: "diggity:distro:id", Value: "centos"},
				{Name: "diggity:distro:prettyName", Value: "CentOS Linux 8"},
				{Name: "diggity:distro:distributionCodename", Value: ""},
				{Name: "diggity:distro:versionID", Value: "8"},
			},
		}},
	}

	for _, test := range tests {
		_output := addDistroComponent(test._distro)

		if _output.Type != test.expected.Type ||
			_output.Name != test.expected.Name ||
			_output.Description != test.expected.Description {
			t.Errorf("Test Failed: Expected output of %v, received: %v ", test.expected, _output)
		}

		if (len(*_output.ExternalReferences) == 0 && len(*test.expected.ExternalReferences) != 0) ||
			(len(*_output.Properties) == 0 && len(*test.expected.Properties) != 0) {
			t.Error("Test Failed: Slice must not be empty.")
		}

		if len(*_output.ExternalReferences) != len(*test.expected.ExternalReferences) {
			t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(*test.expected.ExternalReferences), len(*_output.ExternalReferences))
		}

		if len(*_output.Properties) != len(*test.expected.Properties) {
			t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(*test.expected.Properties), len(*_output.Properties))
		}

		expectedExRefs := *test.expected.ExternalReferences
		for i, exRef := range *_output.ExternalReferences {
			if exRef != expectedExRefs[i] {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", expectedExRefs[i], exRef)
			}
		}

		expectedProperties := *test.expected.Properties
		for i, p := range *_output.Properties {
			if p != expectedProperties[i] {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", expectedProperties[i], p)
			}
		}
	}
}

func TestGetFromSource(t *testing.T) {
	if time.Now().Format(time.RFC3339) != getFromSource().Timestamp {
		t.Errorf("Test Failed: Must produce current timestamp of %v, received: %v.",
			time.Now().Format(time.RFC3339), getFromSource().Timestamp)
	}

	outputTools := *getFromSource().Tools
	if outputTools[0].Vendor != vendor ||
		outputTools[0].Name != name {
		t.Errorf("Test Failed: Expected output of %v, received: %v ", output.Tool{
			Vendor: vendor,
			Name:   name,
		}, outputTools)
	}
}

func TestConvertToComponent(t *testing.T) {
	tests := []ConvertToComponentResult{
		{&cdxPackage1, output.Component{
			BOMRef:     "pkg:alpine/zlib@1.2.12-r3?arch=x86_64&upstream=zlib&distro=?package-id=",
			Type:       library,
			Name:       "zlib",
			Version:    "1.2.12-r3",
			PackageURL: string(model.PURL("pkg:alpine/zlib@1.2.12-r3?arch=x86_64\u0026upstream=zlib\u0026distro=")),
			Licenses: &[]output.License{
				{ID: "Zlib"},
			},
			Properties: &[]output.Property{
				{Name: "diggity:package:type", Value: "apk"},
				{Name: "diggity:cpe23", Value: "cpe:2.3:a:zlib:zlib:1.2.12-r3:*:*:*:*:*:*:*"},
				{Name: "diggity:location:0:layerHash", Value: "9b7240956cfbfefddcd91a2195bfb2ed2cd17bdff81f21111849d643dfaf8131"},
				{Name: "diggity:location:0:path", Value: filepath.Join("lib", "apk", "db", "installed")},
			},
		}},
		{&cdxPackage2, output.Component{
			BOMRef:     "pkg:deb/libapt-pkg6.0@2.2.4arch=s390x?package-id=",
			Type:       library,
			Name:       "libapt-pkg6.0",
			Version:    "2.2.4",
			PackageURL: string(model.PURL("pkg:deb/libapt-pkg6.0@2.2.4arch=s390x")),
			Licenses: &[]output.License{
				{ID: "GPLv2+"},
			},
			Properties: &[]output.Property{
				{Name: "diggity:package:type", Value: "deb"},
				{Name: "diggity:cpe23", Value: "cpe:2.3:a:libapt-pkg6.0:libapt-pkg6.0:2.2.4:*:*:*:*:*:*:*"},
				{Name: "diggity:location:0:layerHash", Value: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0"},
				{Name: "diggity:location:0:path", Value: filepath.Join("var", "lib", "dpkg", "status")},
				{Name: "diggity:location:1:layerHash", Value: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0"},
				{Name: "diggity:location:1:path", Value: filepath.Join("var", "lib", "dpkg", "status")},
			},
		}},
		{&cdxPackage3, output.Component{
			BOMRef:     "pkg:rpm/hardlink@1.0arch=x86_64?package-id=",
			Type:       library,
			Name:       "hardlink",
			Version:    "1.0",
			PackageURL: string(model.PURL("pkg:rpm/hardlink@1.0arch=x86_64")),
			Licenses: &[]output.License{
				{ID: "GPL+"},
			},
			Properties: &[]output.Property{
				{Name: "diggity:package:type", Value: "rpm"},
				{Name: "diggity:cpe23", Value: "cpe:2.3:a:redhat:hardlink:1.0-19.el7:*:*:*:*:*:*:*"},
				{Name: "diggity:cpe23", Value: "cpe:2.3:a:hardlink:hardlink:1.0-19.el7:*:*:*:*:*:*:*"},
				{Name: "diggity:location:0:layerHash", Value: "d1fd2cca7a7751ca9786b088cf639e65088fa0bda34492bb5ba292c32195461a"},
				{Name: "diggity:location:0:path", Value: filepath.Join("var", "lib", "rpm", "Packages")},
			},
		}},
		{&cdxPackage4, output.Component{
			BOMRef:     "pkg:composer/phpdocumentor/reflection@5.2.0?package-id=",
			Type:       library,
			Name:       "phpdocumentor/reflection",
			Version:    "5.2.0",
			PackageURL: string(model.PURL("pkg:composer/phpdocumentor/reflection@5.2.0")),
			Licenses: &[]output.License{
				{ID: "MIT"},
			},
			Properties: &[]output.Property{
				{Name: "diggity:package:type", Value: "php"},
				{Name: "diggity:cpe23", Value: "cpe:2.3:a:phpdocumentor:reflection:5.2.0:*:*:*:*:*:*:*"},
				{Name: "diggity:cpe23", Value: "cpe:2.3:a:reflection:reflection:5.2.0:*:*:*:*:*:*:*"},
				{Name: "diggity:location:0:layerHash", Value: "12a3251e94a5184b3c5f4efbc0c8df91cf8479af3745941c9d9102298d258b83"},
				{Name: "diggity:location:0:path", Value: filepath.Join("opt", "phpdoc", "composer.lock")},
			},
		}},
	}

	for _, test := range tests {
		_output := convertToComponent(test._package)

		if _output.BOMRef != test.expected.BOMRef ||
			_output.Type != test.expected.Type ||
			_output.Name != test.expected.Name ||
			_output.Version != test.expected.Version ||
			_output.PackageURL != test.expected.PackageURL {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, _output)
		}

		expectedProperties := *test.expected.Properties
		expectedLicenses := *test.expected.Licenses

		if (len(*test.expected.Licenses) == 0 && len(*_output.Licenses) != 0) ||
			(len(*test.expected.Properties) == 0 && len(*_output.Properties) != 0) {
			t.Error("Test Failed: Slice must not be empty.")
		}

		if len(*_output.Licenses) != len(*test.expected.Licenses) {
			t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(*test.expected.Licenses), len(*_output.Licenses))
		}

		for i, l := range *_output.Licenses {
			if l != expectedLicenses[i] {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", expectedLicenses[i], l)
			}
		}

		if len(*_output.Properties) != len(*test.expected.Properties) {
			t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(*test.expected.Properties), len(*_output.Properties))
		}

		for i, p := range *_output.Properties {
			if p != expectedProperties[i] {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", expectedProperties[i], p)
			}
		}
	}
}

func TestInitProperties(t *testing.T) {
	tests := []InitPropertiesResult{
		{&cdxPackage1, []output.Property{
			{Name: "diggity:package:type", Value: "apk"},
			{Name: "diggity:cpe23", Value: "cpe:2.3:a:zlib:zlib:1.2.12-r3:*:*:*:*:*:*:*"},
			{Name: "diggity:location:0:layerHash", Value: "9b7240956cfbfefddcd91a2195bfb2ed2cd17bdff81f21111849d643dfaf8131"},
			{Name: "diggity:location:0:path", Value: filepath.Join("lib", "apk", "db", "installed")},
		}},
		{&cdxPackage2, []output.Property{
			{Name: "diggity:package:type", Value: "deb"},
			{Name: "diggity:cpe23", Value: "cpe:2.3:a:libapt-pkg6.0:libapt-pkg6.0:2.2.4:*:*:*:*:*:*:*"},
			{Name: "diggity:location:0:layerHash", Value: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0"},
			{Name: "diggity:location:0:path", Value: filepath.Join("var", "lib", "dpkg", "status")},
			{Name: "diggity:location:1:layerHash", Value: "f1a5f5ce6b163fac7f09b47645c56d2ab676bdcdb268eef06a4d9b782a75bfd0"},
			{Name: "diggity:location:1:path", Value: filepath.Join("var", "lib", "dpkg", "status")},
		}},
		{&cdxPackage3, []output.Property{
			{Name: "diggity:package:type", Value: "rpm"},
			{Name: "diggity:cpe23", Value: "cpe:2.3:a:redhat:hardlink:1.0-19.el7:*:*:*:*:*:*:*"},
			{Name: "diggity:cpe23", Value: "cpe:2.3:a:hardlink:hardlink:1.0-19.el7:*:*:*:*:*:*:*"},
			{Name: "diggity:location:0:layerHash", Value: "d1fd2cca7a7751ca9786b088cf639e65088fa0bda34492bb5ba292c32195461a"},
			{Name: "diggity:location:0:path", Value: filepath.Join("var", "lib", "rpm", "Packages")},
		}},
		{&cdxPackage4, []output.Property{
			{Name: "diggity:package:type", Value: "php"},
			{Name: "diggity:cpe23", Value: "cpe:2.3:a:phpdocumentor:reflection:5.2.0:*:*:*:*:*:*:*"},
			{Name: "diggity:cpe23", Value: "cpe:2.3:a:reflection:reflection:5.2.0:*:*:*:*:*:*:*"},
			{Name: "diggity:location:0:layerHash", Value: "12a3251e94a5184b3c5f4efbc0c8df91cf8479af3745941c9d9102298d258b83"},
			{Name: "diggity:location:0:path", Value: filepath.Join("opt", "phpdoc", "composer.lock")},
		}},
	}

	for _, test := range tests {
		_output := *initProperties(test._package)
		for i, p := range _output {
			if p.Name != test.expected[i].Name || p.Value != test.expected[i].Value {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected[i], p)
			}
		}
	}

}

func TestAddID(t *testing.T) {
	tests := []AddIDResult{
		{&model.Package{
			ID:   "test-id",
			PURL: "test-PURL",
		},
			"test-PURL?package-id=test-id",
		},
		{&model.Package{
			ID:   "123",
			PURL: "456",
		},
			"456?package-id=123",
		},
		{&model.Package{
			ID:   "",
			PURL: "",
		},
			"?package-id=",
		},
		{&model.Package{},
			"?package-id=",
		},
		{&model.Package{
			ID:   "50917e31-97a5-4503-ae5f-789c8e0dca45",
			PURL: "pkg:alpine/ssl_client@1.35.0-r17?arch=x86_64&amp;upstream=busybox&amp;distro=",
		},
			"pkg:alpine/ssl_client@1.35.0-r17?arch=x86_64&amp;upstream=busybox&amp;distro=?package-id=50917e31-97a5-4503-ae5f-789c8e0dca45",
		},
		{&model.Package{
			ID:   "ca0220df-b2e7-4d24-985f-91640664463f",
			PURL: "pkg:rpm/util-linux@2.32.1arch=x86_64",
		},
			"pkg:rpm/util-linux@2.32.1arch=x86_64?package-id=ca0220df-b2e7-4d24-985f-91640664463f",
		},
		{&model.Package{
			ID:   "36177d9c-284c-4300-a22b-d50db3f59dab",
			PURL: "pkg:deb/libudev1@247.3-7arch=s390x",
		},
			"pkg:deb/libudev1@247.3-7arch=s390x?package-id=36177d9c-284c-4300-a22b-d50db3f59dab",
		},
	}

	for _, test := range tests {
		if _output := addID(test._package); _output != test.expected {
			t.Errorf("Test Failed: Input %v must have output of %v, received: %v", test._package, test.expected, _output)
		}
	}
}

func TestConvertLicense(t *testing.T) {
	tests := []ConvertLicenseResult{
		{&model.Package{
			Licenses: []string{
				"MIT",
			},
		},
			[]output.License{
				{ID: "MIT"},
			},
		},
		{&model.Package{
			Licenses: []string{
				"GPL-2.0-only",
			},
		},
			[]output.License{
				{ID: "GPL-2.0-only"},
			},
		},
		{&model.Package{
			Licenses: []string{
				"MIT",
				"BSD",
				"GPL2+",
			},
		},
			[]output.License{
				{ID: "MIT"},
				{ID: "BSD"},
				{ID: "GPL2+"},
			},
		},
		{&model.Package{
			Licenses: []string{
				"test-1",
				"test-2",
				"test-3",
				"test-4",
				"test-5",
			},
		},
			[]output.License{
				{ID: "test-1"},
				{ID: "test-2"},
				{ID: "test-3"},
				{ID: "test-4"},
				{ID: "test-5"},
			},
		},
		{&model.Package{
			Licenses: []string{},
		},
			nil,
		},
	}

	for _, test := range tests {
		_output := convertLicense(test._package)

		if test.expected == nil {
			if _output != nil {
				t.Error("Test Failed: Expected output of nil.")
			}
			continue
		}

		if len(*_output) != len(test.expected) {
			t.Errorf("Test Failed: Expected length of %+v, Received: %+v.", len(test.expected), len(*_output))
		}

		for i, license := range *_output {
			if license.ID != test.expected[i].ID {
				t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected[i], license.ID)
			}
		}
	}
}
