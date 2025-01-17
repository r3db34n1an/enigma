package settings

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"slices"
	"strings"
)

type ExportSetting struct {
	Rotors    []ExportRotor   `json:"rotors,omitempty"     yaml:"rotors,omitempty"`     // Walzenlage
	Reflector string          `json:"reflector,omitempty"  yaml:"reflector,omitempty"`  // Reflektor
	PlugBoard ExportPlugBoard `json:"plug_board,omitempty" yaml:"plug_board,omitempty"` // Steckerverbindungen

	// generated values
	RotorInfo     []string `json:"rotor_info,omitempty"     yaml:"rotor_info,omitempty"`
	RotorSettings string   `json:"rotor_settings,omitempty" yaml:"rotor_settings,omitempty"`
	Plugs         string   `json:"plugs,omitempty"          yaml:"plugs,omitempty"`
	Key           string   `json:"key,omitempty"            yaml:"key,omitempty"`
}

func (what *ExportSetting) Parse(value string) error {
	parseError := yaml.Unmarshal([]byte(value), &what)
	if parseError != nil {
		return fmt.Errorf("failed to parse setting: %v", parseError)
	}

	return nil
}

func (what *ExportSetting) Print() (string, error) {
	value, marshalError := yaml.Marshal(what)
	if marshalError != nil {
		return "", fmt.Errorf("failed to marshal setting: %v", marshalError)
	}

	return string(value), nil
}

func (what *ExportSetting) Generate() error {
	for index := range what.Rotors {
		what.RotorInfo = append(what.RotorInfo, fmt.Sprintf("%v: %v + %v", what.Rotors[index].Name, what.Rotors[index].Position, what.Rotors[index].RingSetting))
	}

	var keys []string
	for key, value := range what.PlugBoard {
		if key != value {
			keys = append(keys, key)
		}
	}

	slices.Sort(keys)
	var plugs []string
	for _, key := range keys {
		plugs = append(plugs, fmt.Sprintf("%v%v", key, what.PlugBoard[key]))
	}

	what.Plugs = strings.Join(plugs, " ")
	what.RotorSettings = strings.Join(what.RotorInfo, "; ")
	what.Key = strings.Join([]string{what.Reflector, what.RotorSettings, what.Plugs}, " | ")

	value, printError := what.Print()
	if printError != nil {
		return fmt.Errorf("failed to print setting: %v", printError)
	}

	what.Key = value
	return nil
}
