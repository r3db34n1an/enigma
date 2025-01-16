package enigma

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

var reflectors Reflectors

type Reflector struct {
	Name    string
	Mapping map[int]int
}

type Reflectors map[string]*Reflector

func GetReflector(name string) (*Reflector, error) {
	if reflectors == nil {
		reflectors = make(Reflectors)
		loadError := reflectors.load(reflectorsYaml)
		if loadError != nil {
			return nil, fmt.Errorf("failed to load reflectors: %v", loadError)
		}
	}

	reflector, ok := reflectors[strings.ToUpper(name)]
	if !ok {
		return nil, fmt.Errorf("reflector %q not found", name)
	}

	return reflector, nil
}

func (what *Reflector) Reflect(in int) int {
	out, ok := what.Mapping[in]
	if !ok {
		return -1
	}

	return out
}

func (what *Reflector) load(data any) error {
	reflectors = make(Reflectors)
	switch castData := data.(type) {
	case string:
		what.Mapping = make(map[int]int)
		for k, v := range strings.ToUpper(castData) {
			what.Mapping[k] = strings.IndexRune(upperCase, v)
			if what.Mapping[k] == -1 {
				return fmt.Errorf("invalid reflector value %q", v)
			}
		}
	}

	return nil
}

func (what *Reflectors) load(data any) error {
	castData, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("invalid reflectors data, expected []byte, got %T", data)
	}

	var items map[string]any
	parseError := yaml.Unmarshal(castData, &items)
	if parseError != nil {
		return fmt.Errorf("failed to parse reflectors: %v", parseError)
	}

	newReflectors := make(map[string]*Reflector)
	for reflectorName, reflectorValue := range items {
		reflector := &Reflector{
			Name: strings.ToUpper(reflectorName),
		}

		loadError := reflector.load(reflectorValue)
		if loadError != nil {
			return fmt.Errorf("failed to load rotor %q: %v", reflectorName, loadError)
		}

		newReflectors[reflector.Name] = reflector
	}

	*what = newReflectors

	return nil
}
