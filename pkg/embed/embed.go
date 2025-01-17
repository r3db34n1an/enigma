package embed

import _ "embed"

//go:embed config/settings.yaml
var SettingsYaml []byte

//go:embed config/rotors.yaml
var RotorsYaml []byte

//go:embed config/reflectors.yaml
var ReflectorsYaml []byte
