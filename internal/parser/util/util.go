package util

import (
	"os"
	"strings"

	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/parser/bom"
)

// ParserNames slice of supported parser names
var ParserNames = []string{
	"apk",
	"debian",
	"java",
	"npm",
	"composer",
	"python",
	"gem",
	"rpm",
	"dart",
	"nuget",
	"go",
	"rust",
	"conan",
	"hackage",
	"pod",
	"hex",
	"portage",
}

// TrimUntilLayer Returns file path without layer hash
func TrimUntilLayer(location model.Location) string {
	directories := strings.Split(location.Path, string(os.PathSeparator))
	index := IndexOf(directories, location.LayerHash) + 1
	directory := ""
	for index < len(directories) {
		if index == len(directories)-1 {
			directory += directories[index]
		} else {
			directory += directories[index] + string(os.PathSeparator)
		}

		index++
	}
	return directory
}

// ParserEnabled checks if all or a specific parser is enabled
func ParserEnabled(parser string) bool {
	if len(*bom.Arguments.EnabledParsers) == 0 {
		return true
	}
	if StringSliceContains(*bom.Arguments.EnabledParsers, parser) {
		return true
	}
	return false
}

// IndexOf returns index of a string from a slice
func IndexOf(array []string, s string) int {
	for idx, a := range array {
		if s == a {
			return idx
		}
	}
	return -1
}

// StringSliceContains checks if a string slice contains specified string
func StringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// SourceIsDir returns if source is dir
func SourceIsDir() bool {
	return len(*bom.Arguments.Dir) > 0
}

// FormatLockKeyVal formats .lock Key Value Data String
func FormatLockKeyVal(kv string) string {
	trimmed := strings.TrimSpace(kv)
	return strings.Replace(trimmed, `"`, "", -1)
}
