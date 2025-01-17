package settings

type ExportRotor struct {
	Name        string `json:"name,omitempty"         yaml:"name,omitempty"`
	Position    string `json:"position,omitempty"     yaml:"position,omitempty"`     // Grundstellung
	RingSetting string `json:"ring_setting,omitempty" yaml:"ring_setting,omitempty"` // Ringstellung
}
