package settings

import (
	"fmt"
	"github.com/r3db34n1an/enigma/pkg/defs"
	"github.com/r3db34n1an/enigma/pkg/embed"
	"gopkg.in/yaml.v3"
	"strings"
)

var rotors Rotors

type Rotor struct {
	Name        string
	Notches     []int // Kerben
	Position    int   // Grundstellung
	RingSetting int   // Ringstellung
	Forward     map[int]int
	Reverse     map[int]int
}

type Rotors map[string]Rotor
type RotorGroup []*Rotor

func GetRotor(name string) (*Rotor, error) {
	if rotors == nil {
		rotors = make(Rotors)
		loadError := rotors.load(embed.RotorsYaml)
		if loadError != nil {
			return nil, fmt.Errorf("failed to load rotors: %v", loadError)
		}
	}

	rotor, ok := rotors[strings.ToUpper(name)]
	if !ok {
		return nil, fmt.Errorf("rotor %q not found", name)
	}

	newRotor := rotor
	return &newRotor, nil
}

func (what *Rotor) ParseRingSetting(name string) error {
	if name == "" {
		return fmt.Errorf("missing ring setting")
	}

	index := strings.IndexRune(defs.UpperCase, rune(strings.ToUpper(name)[0]))
	if index < 0 {
		return fmt.Errorf("invalid ring setting %q", name)
	}

	what.RingSetting = index + 1
	return nil
}

func (what *Rotor) ParsePosition(name string) error {
	if name == "" {
		return fmt.Errorf("missing ring setting")
	}

	index := strings.IndexRune(defs.UpperCase, rune(strings.ToUpper(name)[0]))
	if index < 0 {
		return fmt.Errorf("invalid ring setting %q", name)
	}

	what.Position = index + 1
	return nil
}

func (what *Rotor) encrypt(in int) int {
	limit := len(defs.UpperCase)

	advance := what.Position - what.RingSetting
	in += limit + advance
	in %= limit
	in = what.Forward[in]
	in += limit - advance
	in %= limit

	return in
}

func (what *Rotor) decrypt(in int) int {
	limit := len(defs.UpperCase)

	in = (in - what.RingSetting + what.Position + limit) % limit
	in = what.Reverse[in]
	in = (in + what.RingSetting - what.Position + limit) % limit

	return in
}

func (what *Rotor) move() {
	limit := len(defs.UpperCase)
	what.Position = (what.Position + 1) % limit
}

func (what *Rotor) shouldMove() bool {
	for _, notch := range what.Notches {
		if what.Position == notch {
			return true
		}
	}

	return false
}

func (what *Rotor) load(rotorValue any) error {
	switch castRotorValue := rotorValue.(type) {
	case map[string]any:
		for rotorAttributeName, rotorAttributeValue := range castRotorValue {
			switch strings.ToLower(rotorAttributeName) {
			case "mapping":
				switch castRotorAttributeValue := rotorAttributeValue.(type) {
				case string:
					if len(castRotorAttributeValue) != 26 {
						return fmt.Errorf("invalid rotor mapping %q, expected 26 characters", rotorAttributeValue)
					}

					what.Forward = make(map[int]int)
					what.Reverse = make(map[int]int)
					for index, value := range strings.ToUpper(castRotorAttributeValue) {
						mapped := strings.IndexRune(defs.UpperCase, value)
						if mapped == -1 {
							return fmt.Errorf("invalid rotor mapping value %q", value)
						}

						_, forwardExists := what.Forward[index]
						if forwardExists {
							return fmt.Errorf("duplicate forward rotor mapping value %q", value)
						}

						_, reverseExists := what.Reverse[mapped]
						if reverseExists {
							return fmt.Errorf("duplicate reverse rotor mapping value %q", value)
						}

						what.Forward[index] = mapped
						what.Reverse[mapped] = index
					}
				}

			case "notches":
				switch castRotorAttributeValue := rotorAttributeValue.(type) {
				case []any:
					what.Notches = make([]int, 0)
					for _, notch := range castRotorAttributeValue {
						switch castNotch := notch.(type) {
						case string:
							if len(castNotch) != 1 {
								return fmt.Errorf("invalid rotor notch %q, expected 1 character", notch)
							}

							what.Notches = append(what.Notches, strings.IndexRune(defs.UpperCase, rune(castNotch[0])))

						default:
							return fmt.Errorf("invalid rotor notch %T, expected int or string", notch)
						}
					}
				}

			default:
				return fmt.Errorf("invalid rotor attiribute %q", rotorAttributeName)
			}
		}

	default:
		return fmt.Errorf("invalid reflector %T, expected map[string]any", rotorValue)
	}

	return nil
}

func (what *Rotors) load(data any) error {
	var items map[string]any
	parseError := yaml.Unmarshal(data.([]byte), &items)
	if parseError != nil {
		return fmt.Errorf("failed to parse rotors: %v", parseError)
	}

	for rotorName, rotorValue := range items {
		rotor := Rotor{
			Name: strings.ToUpper(rotorName),
		}

		loadError := rotor.load(rotorValue)
		if loadError != nil {
			return fmt.Errorf("failed to load rotor %q: %v", rotorName, loadError)
		}

		rotors[rotor.Name] = rotor
	}

	return nil
}

func (what *RotorGroup) Encrypt(in int) int {
	return what.encrypt(in, *what)
}

func (what *RotorGroup) Decrypt(in int) int {
	return what.decrypt(in, *what)
}

func (what *RotorGroup) Move() {
	if len(*what) < 3 {
		return
	}

	rightIndex := len(*what) - 1
	middleIndex := rightIndex - 1
	leftIndex := middleIndex - 1

	middleShouldMove := (*what)[rightIndex].shouldMove() || (*what)[middleIndex].shouldMove()
	leftShouldMove := (*what)[middleIndex].shouldMove()

	(*what)[rightIndex].move()

	if middleShouldMove {
		(*what)[middleIndex].move()
	}

	if leftShouldMove {
		(*what)[leftIndex].move()
	}
}

func (what *RotorGroup) encrypt(in int, rotors RotorGroup) int {
	if len(rotors) == 0 {
		return in
	}

	in = what.encrypt(in, rotors[1:])
	in = rotors[0].encrypt(in)

	return in
}

func (what *RotorGroup) decrypt(in int, rotors RotorGroup) int {
	if len(rotors) == 0 {
		return in
	}

	in = rotors[0].decrypt(in)
	in = what.decrypt(in, rotors[1:])

	return in
}
