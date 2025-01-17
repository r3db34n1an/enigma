package settings

import (
	"fmt"
	"github.com/r3db34n1an/enigma/pkg/defs"
	"github.com/r3db34n1an/enigma/pkg/embed"
	"gopkg.in/yaml.v3"
	"strings"
)

var settings Settings

type Setting struct {
	IDGroups  []string   // Kenngruppen
	Rotors    RotorGroup // Walzenlage
	Reflector Reflector  // Reflektor
	PlugBoard PlugBoard  // Steckerverbindungen
}

type Settings []Setting

func (what *Setting) Get(name string) error {
	if settings == nil {
		loadError := settings.Load(embed.SettingsYaml)
		if loadError != nil {
			return fmt.Errorf("failed to load settings: %v", loadError)
		}
	}

	for _, setting := range settings {
		for _, idGroup := range setting.IDGroups {
			if strings.ToUpper(name) == idGroup {
				*what = setting
				return nil
			}
		}
	}

	return fmt.Errorf("setting %q not found", name)
}

func (what *Setting) Random() error {
	if settings == nil {
		loadError := settings.Load(embed.SettingsYaml)
		if loadError != nil {
			return fmt.Errorf("failed to load settings: %v", loadError)
		}
	}

	*what = settings[defs.RandomInt(0, len(settings)-1)]
	for _, rotor := range what.Rotors {
		rotor.Position = defs.RandomInt(0, 25)
	}

	return nil
}

func (what *Setting) Export() ExportSetting {
	var exportedRotors []ExportRotor
	for _, rotor := range what.Rotors {
		exportedRotors = append(exportedRotors, ExportRotor{
			Name:        rotor.Name,
			Position:    string(rune(rotor.Position + 'A')),
			RingSetting: string(rune(rotor.RingSetting + 'A')),
		})
	}

	exportedPlugBoard := make(map[string]string)
	for plug, value := range what.PlugBoard.Mapping {
		exportedPlugBoard[string(defs.UpperCase[plug])] = string(defs.UpperCase[value])
	}

	return ExportSetting{
		Rotors:    exportedRotors,
		Reflector: what.Reflector.Name,
		PlugBoard: exportedPlugBoard,
	}
}

func (what *Setting) Import(exportSetting ExportSetting) error {
	what.Rotors = nil
	for _, exportedRotor := range exportSetting.Rotors {
		rotorError := what.ImportRotor(exportedRotor)
		if rotorError != nil {
			return rotorError
		}
	}

	reflectorError := what.ImportReflector(exportSetting.Reflector)
	if reflectorError != nil {
		return reflectorError
	}

	plugBoardError := what.ImportPlugBoard(exportSetting.PlugBoard)
	if plugBoardError != nil {
		return plugBoardError
	}

	return what.validate(true)
}

func (what *Setting) ImportRotor(exportRotor ExportRotor) error {
	rotor, rotorError := GetRotor(strings.ToUpper(exportRotor.Name))
	if rotorError != nil {
		return fmt.Errorf("invalid rotor %q: %v", exportRotor.Name, rotorError)
	}

	exportRotor.Position = strings.ToUpper(exportRotor.Position)
	exportRotor.RingSetting = strings.ToUpper(exportRotor.RingSetting)

	rotor.Position = strings.IndexRune(defs.UpperCase, rune(exportRotor.Position[0]))
	rotor.RingSetting = strings.IndexRune(defs.UpperCase, rune(exportRotor.RingSetting[0]))

	what.Rotors = append(what.Rotors, rotor)

	return nil
}

func (what *Setting) ImportReflector(exportReflector string) error {
	reflector, reflectorError := GetReflector(exportReflector)
	if reflectorError != nil {
		return fmt.Errorf("invalid reflector %q: %v", exportReflector, reflectorError)
	}

	what.Reflector = *reflector

	return nil
}

func (what *Setting) ImportPlugBoard(exportPlugBoard ExportPlugBoard) error {
	what.PlugBoard.Mapping = make(map[int]int)

	if exportPlugBoard == nil {
		return nil
	}

	for plug, value := range exportPlugBoard {
		plugIndex := strings.IndexRune(defs.UpperCase, rune(plug[0]))
		valueIndex := strings.IndexRune(defs.UpperCase, rune(value[0]))

		what.PlugBoard.Mapping[plugIndex] = valueIndex
		what.PlugBoard.Mapping[valueIndex] = plugIndex
	}

	return nil
}

func (what *Setting) Clone() (*Setting, error) {
	var setting Setting
	importError := setting.Import(what.Export())
	if importError != nil {
		return nil, fmt.Errorf("failed to clone setting: %v", importError)
	}

	return &setting, nil
}

func (what *Setting) Load(data any) error {
	switch castData := data.(type) {
	case map[string]any:
		for key, value := range castData {
			switch strings.ToLower(key) {
			case "id_groups":
				importError := what.LoadIDGroups(value)
				if importError != nil {
					return fmt.Errorf("invalid id_groups: %v", importError)
				}

			case "rotors":
				importError := what.LoadRotors(value)
				if importError != nil {
					return fmt.Errorf("invalid rotors: %v", importError)
				}

			case "reflector":
				importError := what.LoadReflector(value)
				if importError != nil {
					return fmt.Errorf("invalid reflector: %v", importError)
				}

			case "plug_board":
				importError := what.LoadPlugBoard(value)
				if importError != nil {
					return fmt.Errorf("invalid plug_board: %v", importError)
				}

			default:
				return fmt.Errorf("invalid setting key %q", key)
			}
		}

	default:
		return fmt.Errorf("invalid setting format %T, expected map[string]any", data)
	}

	return what.validate(false)
}

func (what *Setting) LoadIDGroups(value any) error {
	switch castValue := value.(type) {
	case []any:
		for _, idGroup := range castValue {
			switch castIDGroup := idGroup.(type) {
			case string:
				if len(castIDGroup) != 3 {
					return fmt.Errorf("invalid id_group %q, expected 3 characters", castIDGroup)
				}

				what.IDGroups = append(what.IDGroups, castIDGroup)

			default:
				return fmt.Errorf("invalid id_groups format %T, expected string", idGroup)
			}
		}

		if len(what.IDGroups) != 4 {
			return fmt.Errorf("invalid id_groups length %d, expected 4", len(what.IDGroups))
		}

	default:
		return fmt.Errorf("invalid id_groups format %T, expected []any", value)
	}

	return nil
}

func (what *Setting) LoadRotors(value any) error {
	switch castValue := value.(type) {
	case []any:
		for _, rotorData := range castValue {
			switch castRotorMap := rotorData.(type) {
			case map[string]any:
				if len(castRotorMap) != 1 {
					return fmt.Errorf("invalid rotor %q, expected single item", castRotorMap)
				}

				for rotorName, rotorValue := range castRotorMap {
					rotor, rotorError := GetRotor(rotorName)
					if rotorError != nil {
						return fmt.Errorf("invalid rotor %q: %v", rotorName, rotorError)
					}

					switch castRotorValue := rotorValue.(type) {
					case int:
						rotor.RingSetting = castRotorValue - 1

					default:
						return fmt.Errorf("invalid rotor ring %T, expected int or string", rotorValue)
					}

					what.Rotors = append(what.Rotors, rotor)
				}

			default:
				return fmt.Errorf("invalid rotor format %T, expected map[string]any", rotorData)
			}
		}

	default:
		return fmt.Errorf("invalid rotors format %T, expected []any", value)
	}

	return nil
}

func (what *Setting) LoadReflector(value any) error {
	switch castValue := value.(type) {
	case string:
		reflector, reflectorError := GetReflector(castValue)
		if reflectorError != nil {
			return fmt.Errorf("invalid reflector %q: %v", castValue, reflectorError)
		}

		what.Reflector = *reflector

	default:
		return fmt.Errorf("invalid reflector %T, expected string", value)
	}

	return nil
}

func (what *Setting) LoadPlugBoard(value any) error {
	what.PlugBoard.Mapping = make(map[int]int)

	switch castValue := value.(type) {
	case string:
		parseError := what.PlugBoard.Parse(castValue)
		if parseError != nil {
			return fmt.Errorf("invalid plug board: %v", parseError)
		}

	default:
		return fmt.Errorf("invalid plug board %T, expected []string", value)
	}

	return nil
}

func (what *Setting) validate(imported bool) error {
	// Validate IDGroups
	if !imported {
		if len(what.IDGroups) != 4 {
			return fmt.Errorf("invalid id_groups length %d, expected 4", len(what.IDGroups))
		}

		for _, idGroup := range what.IDGroups {
			if len(idGroup) != 3 {
				return fmt.Errorf("invalid id_group %q, expected 3 characters", idGroup)
			}
		}
	}

	// Validate Rotors
	if len(what.Rotors) < 3 || len(what.Rotors) > 4 {
		return fmt.Errorf("invalid rotors length %d, expected 3 or 4", len(what.Rotors))
	}

	for _, rotor := range what.Rotors {
		if rotor == nil {
			return fmt.Errorf("invalid rotor: nil")
		}

		if rotor.RingSetting < 0 || rotor.RingSetting > 25 {
			return fmt.Errorf("invalid rotor ring %d, expected 0-25", rotor.RingSetting)
		}

		if rotor.Forward == nil {
			return fmt.Errorf("invalid rotor mapping: nil")
		}

		if len(rotor.Forward) != 26 {
			return fmt.Errorf("invalid rotor mapping length %d, expected 26", len(rotor.Forward))
		}

		for index, value := range rotor.Forward {
			if value < 0 || value > 25 {
				return fmt.Errorf("invalid rotor mapping value %d, expected 0-25", value)
			}

			if _, exists := rotor.Forward[index]; !exists {
				return fmt.Errorf("invalid rotor mapping index %d, expected 0-25", index)
			}
		}

		if rotor.Reverse == nil {
			return fmt.Errorf("invalid rotor mapping: nil")
		}

		if len(rotor.Reverse) != 26 {
			return fmt.Errorf("invalid rotor mapping length %d, expected 26", len(rotor.Reverse))
		}

		for index, value := range rotor.Reverse {
			if value < 0 || value > 25 {
				return fmt.Errorf("invalid rotor mapping value %d, expected 0-25", value)
			}

			if _, exists := rotor.Reverse[index]; !exists {
				return fmt.Errorf("invalid rotor mapping index %d, expected 0-25", index)
			}
		}
	}

	// Validate Reflector
	if len(what.Reflector.Mapping) != 26 {
		return fmt.Errorf("invalid reflector mapping length %d, expected 26", len(what.Reflector.Mapping))
	}

	for index, value := range what.Reflector.Mapping {
		if value < 0 || value > 25 {
			return fmt.Errorf("invalid reflector mapping value %d, expected 0-25", value)
		}

		if _, exists := (what.Reflector.Mapping)[index]; !exists {
			return fmt.Errorf("invalid reflector mapping index %d, expected 0-25", index)
		}
	}

	// Validate PlugBoard
	for index, value := range what.PlugBoard.Mapping {
		if value < 0 || value > 25 {
			return fmt.Errorf("invalid plug board value %d, expected 0-25", value)
		}

		if _, exists := what.PlugBoard.Mapping[index]; !exists {
			return fmt.Errorf("invalid plug board forward index %d, expected 0-25", index)
		}
	}

	return nil
}

func (what *Settings) Load(data any) error {
	var items []any
	parseError := yaml.Unmarshal(data.([]byte), &items)
	if parseError != nil {
		return fmt.Errorf("failed to parse settings: %v", parseError)
	}

	for _, settingData := range items {
		setting := new(Setting)
		settingError := setting.Load(settingData)
		if settingError != nil {
			return settingError
		}

		*what = append(*what, *setting)
	}

	return nil
}
