package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"

	"github.com/alessio/shellescape"
	"github.com/tidwall/jsonc"
)

type Settings struct {
	CMakeConfigureSettings map[string]string `json:"cmake.configureSettings"`
}

func ReadSettings(input string) (Settings, error) {
	s, err := os.ReadFile(input)

	if err != nil {
		return Settings{}, err
	}

	var settings Settings
	err = json.Unmarshal(jsonc.ToJSON(s), &settings)
	if err != nil {
		return Settings{}, err
	}
	return settings, nil
}

func (s Settings) ExtractCMakeConfigureSettingCommandline() []string {
	var args []string
	for key, value := range s.CMakeConfigureSettings {
		//fmt.Println(key, value)
		args = append(args, fmt.Sprintf("-D%s=%v", key, shellescape.Quote(value)))
	}
	sort.Strings(args)
	return args
}

func ShowCommand(args []string) {
	fmt.Printf("%v\n", args)

	cmd := exec.Command("cmake", args...)
	res := cmd.Run()
	if res != nil {
		fmt.Printf("error: %v\n", res)
	}
	cmd.Wait()
	os.Exit(cmd.ProcessState.ExitCode())
}

func GetEnvOrDefault(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

func main() {
	//inFile := flag.String("i", ".vscode/settings.json", "input file, if not $PWD/.vscode/settings.json")
	//flag.Parse()

	inFile := GetEnvOrDefault("VSCODE_SETTINGS", ".vscode/settings.json")

	settings, err := ReadSettings(inFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	args := settings.ExtractCMakeConfigureSettingCommandline()

	//fmt.Print(strings.Join(args, " "))

	args = append(args, os.Args[1:]...)

	cmd := exec.Command("cmake", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Println(cmd.Args)

	res := cmd.Run()

	if res != nil {
		fmt.Printf("error: %v\n", res)
	}

	os.Exit(cmd.ProcessState.ExitCode())
}
