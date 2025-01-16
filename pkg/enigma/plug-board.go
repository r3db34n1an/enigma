package enigma

import (
	"fmt"
	"strings"
)

type PlugBoard struct {
	Mapping map[int]int
}

func (what *PlugBoard) Transform(in int) int {
	out, ok := what.Mapping[in]
	if !ok {
		return in
	}

	return out
}

func (what *PlugBoard) Parse(in string) error {
	for _, plug := range strings.Fields(strings.ToUpper(in)) {
		if len(plug) != 2 {
			return fmt.Errorf("invalid plug board value %q, expected 2 characters", plug)
		}

		plugOne := strings.IndexRune(upperCase, rune(plug[0]))
		plugTwo := strings.IndexRune(upperCase, rune(plug[1]))
		if plugOne == -1 || plugTwo == -1 {
			return fmt.Errorf("invalid plug board value %q", plug)
		}

		_, plugOneExists := what.Mapping[plugOne]
		if plugOneExists {
			return fmt.Errorf("duplicate plug board value %q", plug)
		}

		_, plugTwoExists := what.Mapping[plugTwo]
		if plugTwoExists {
			return fmt.Errorf("duplicate plug board value %q", plug)
		}

		what.Mapping[plugOne] = plugTwo
		what.Mapping[plugTwo] = plugOne
	}

	for _, letter := range upperCase {
		plug := strings.IndexRune(upperCase, letter)

		_, forwardExists := what.Mapping[strings.IndexRune(upperCase, letter)]
		if !forwardExists {
			what.Mapping[plug] = plug
		}

		_, reverseExists := what.Mapping[strings.IndexRune(upperCase, letter)]
		if !reverseExists {
			what.Mapping[plug] = plug
		}
	}

	return nil
}
