// Copyright 2022 Markus Holmström (MawKKe) markus@mawkke.fi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"reflect"
	"testing"
)

// Well, technically this is not valid "JSON" but.. whatever
var exampleSettingsJSON = `
{
    "editor.formatOnSave": true,
    "cmake.configureOnOpen": true,
    "cmake.configureArgs": [
        "-GNinja"
    ],
	// A comment here
    "cmake.configureSettings": {
		"CMAKE_CXX_COMPILER": "clang++",
		"CMAKE_CXX_FLAGS_INIT": "-fdiagnostics-color=always -O3",
		"CMAKE_CXX_STANDARD_REQUIRED": "ON", // stupid CMake does not put -std= flag on the command line with GCC, but on Clang it is present
		"CMAKE_CXX_STANDARD": "17"
    },
    "cmake.ctestArgs": []
	// and rest of your settings.json
}
`

func TestParseVSCodeSettings(t *testing.T) {
	settings, err := ParseVSCodeSettings([]byte(exampleSettingsJSON))
	if err != nil {
		t.Fatalf("VSCode settings parsing failed: %q", err)
	}

	t.Run("Test that cmake.configureArgs is parsed correctly", func(t *testing.T) {
		expectedArguments := []string{"-GNinja"}
		if !reflect.DeepEqual(settings.CMakeConfigureArguments, expectedArguments) {
			t.Fatalf("Expected CMakeConfigureArguments: %v, got: %v",
				expectedArguments, settings.CMakeConfigureArguments)
		}
	})

	t.Run("Test that cmake.configureSettings is parsed correctly", func(t *testing.T) {
		expectedSettings := map[string]string{
			"CMAKE_CXX_COMPILER":          "clang++",
			"CMAKE_CXX_FLAGS_INIT":        "-fdiagnostics-color=always -O3",
			"CMAKE_CXX_STANDARD":          "17",
			"CMAKE_CXX_STANDARD_REQUIRED": "ON",
		}

		if !reflect.DeepEqual(expectedSettings, settings.CMakeConfigureSettings) {
			t.Fatalf("Expected CMakeConfigureSettings:\n\t%v, got:\n\t%v",
				settings.CMakeConfigureSettings, expectedSettings)
		}
	})

	t.Run("Test that cmake.configureSettings are formatted properly", func(t *testing.T) {
		formatted := settings.FormatCMakeConfigureSettings()
		// NOTE: The program sorts the keys since maps are iterated in random order
		expected := []string{
			"-DCMAKE_CXX_COMPILER=clang++",
			"-DCMAKE_CXX_FLAGS_INIT='-fdiagnostics-color=always -O3'",
			"-DCMAKE_CXX_STANDARD=17",
			"-DCMAKE_CXX_STANDARD_REQUIRED=ON",
		}
		if !reflect.DeepEqual(formatted, expected) {
			t.Fatalf("Expected formatted:\n\t%v,\ngot:\n\t%v", expected, formatted)
		}
	})
	t.Run("Test that computed CLI arguments are correct", func(t *testing.T) {
		cliArgs := settings.CollectCLIArgs("-h")
		expected := []string{
			"-DCMAKE_CXX_COMPILER=clang++",
			"-DCMAKE_CXX_FLAGS_INIT='-fdiagnostics-color=always -O3'",
			"-DCMAKE_CXX_STANDARD=17",
			"-DCMAKE_CXX_STANDARD_REQUIRED=ON",
			"-GNinja",
			"-h",
		}

		if !reflect.DeepEqual(expected, cliArgs) {
			t.Fatalf("Expected command line:\n\t%v,\ngot:\n\t%v", expected, cliArgs)
		}
	})
}
