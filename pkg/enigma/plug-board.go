package enigma

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
