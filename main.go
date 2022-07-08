package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/tidwall/jsonc"
)

// VSCodeSettings is a struct representing VCode settings.json relating to CMake options
type VSCodeSettings struct {
	CMakeConfigureSettings  map[string]string `json:"cmake.configureSettings"`
	CMakeConfigureArguments []string          `json:"cmake.configureArgs"`
}

// ReadVSCodeSettings extracts CMake -DKEY=VALUE parameters from given input file
func ReadVSCodeSettings(inputFile string) (VSCodeSettings, error) {
	contents, err := os.ReadFile(inputFile)
	if err != nil {
		return VSCodeSettings{}, err
	}
	return ParseVSCodeSettings(contents)
}

// ParseVSCodeSettings extracts CMake -DKEY=VALUE parameters from given input byte slice
func ParseVSCodeSettings(inputString []byte) (VSCodeSettings, error) {
	var settings VSCodeSettings
	// We can't do normal JSON decode, since the file may contain
	// comments (which makes it non-standard/invalid JSON). We use 'jsonc' library
	// for transforming the input into suitable, valid JSON.
	err := json.Unmarshal(jsonc.ToJSON(inputString), &settings)
	if err != nil {
		return VSCodeSettings{}, err
	}
	return settings, nil
}

// FormatCMakeConfigureSettings produces a list of "-DKEY=VALUE" arguments
// from the configure settings, suitable for passing to CMake program.
func (settings VSCodeSettings) FormatCMakeConfigureSettings() []string {
	var args []string
	for key, value := range settings.CMakeConfigureSettings {
		//fmt.Println(key, value)
		args = append(args, fmt.Sprintf("-D%s=%s", key, shellescape.Quote(value)))
	}
	// golang iterates map items in random order; this should ensure deterministic results.
	sort.Strings(args)
	return args
}

// CollectCLIArgs builds a complete set of CMake command line arguments from
// all known information.
func (settings VSCodeSettings) CollectCLIArgs(argv ...string) []string {

	var allArgs []string
	allArgs = append(allArgs, settings.FormatCMakeConfigureSettings()...)
	allArgs = append(allArgs, settings.CMakeConfigureArguments...)
	allArgs = append(allArgs, argv...)
	return allArgs
}

// RunCMakeConfigure run CMake configuration command using the given settings.
func RunCMakeConfigure(settings VSCodeSettings, dryRun bool) int {

	cmd := exec.Command("cmake", settings.CollectCLIArgs(os.Args[1:]...)...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("Running command:\n\t%v\n\n", strings.Join(cmd.Args, " "))

	if dryRun {
		return 0
	}

	if res := cmd.Run(); res != nil {
		fmt.Printf("error: %v\n", res)
	}

	return cmd.ProcessState.ExitCode()
}

// GetEnvOrDefault returns environment variable described by 'key', or fallback
// if the given key does not exist (or is empty).
func GetEnvOrDefault(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

var helpText = `==========

%[1]s:
	This a wrapper program for configuring CMake project automatically using contents of .vscode/settings.json

	Most of the time you'll call it once for configuring the project, and then resume with normal CMake:

		$ %[1]s -B mybuild .
		$ cmake --build mybuild

	You may also perform a dry-run by enabling the VCC_DRY_RUN environment value (with a non-"FALSE" value):

		$ env VCC_DRY_RUN=1 %[1]s -B mybuild .

	...and then run the shown command manually in a terminal.

	The program assumes the VSCode settings file to be found under $PWD/.vscode/settings.json,
	but you may specify alternative path via

		$ env VCC_VSCODE_SETTINGS=path/to/mysettings.json %[1]s -B mybuild .

	Of course, a combination of these environment variables should work as expected.

==========

`

func showHelp() {
	fmt.Printf(helpText, os.Args[0])
}

func main() {
	inFile := GetEnvOrDefault("VCC_VSCODE_SETTINGS", ".vscode/settings.json")
	dryRun := GetEnvOrDefault("VCC_DRY_RUN", "FALSE") != "FALSE"

	if len(os.Args) >= 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		showHelp()
	}

	settings, err := ReadVSCodeSettings(inFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	retcode := RunCMakeConfigure(settings, dryRun)

	os.Exit(retcode)
}
