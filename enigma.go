package enigma

import (
	"fmt"
	"github.com/r3db34n1an/enigma/pkg/enigma"
	"github.com/r3db34n1an/enigma/pkg/settings"
)

type Enigma struct {
	machine *enigma.Enigma
}

func NewEnigma(copyExtra bool, preserveCase bool) (*Enigma, error) {
	machine, enigmaError := enigma.NewEnigma(copyExtra, preserveCase)
	if enigmaError != nil {
		return nil, fmt.Errorf("failed to create enigma machine: %v", enigmaError)
	}

	return &Enigma{
		machine: machine,
	}, nil
}

func (what *Enigma) Init() error {
	if what.machine == nil {
		return fmt.Errorf("no enigma machine")
	}

	return what.machine.Init()
}

func (what *Enigma) Encrypt(plainText []byte, key string) ([]byte, error) {
	if what.machine == nil {
		return nil, fmt.Errorf("no enigma machine")
	}

	return what.machine.Encrypt(plainText, key)
}

func (what *Enigma) EncryptWithPlugBoard(plainText []byte, key string, plugBoard string) ([]byte, error) {
	if what.machine == nil {
		return nil, fmt.Errorf("no enigma machine")
	}

	return what.machine.EncryptWithPlugBoard(plainText, key, plugBoard)
}

func (what *Enigma) EncryptWithSetting(plainText []byte, setting *settings.Setting) ([]byte, error) {
	if what.machine == nil {
		return nil, fmt.Errorf("no enigma machine")
	}

	return what.machine.EncryptWithSetting(plainText, setting)
}

func (what *Enigma) Decrypt(cipherText []byte, key string) ([]byte, error) {
	if what.machine == nil {
		return nil, fmt.Errorf("no enigma machine")
	}

	return what.machine.Decrypt(cipherText, key)
}

func (what *Enigma) DecryptWithPlugBoard(cipherText []byte, key string, plugBoard string) ([]byte, error) {
	if what.machine == nil {
		return nil, fmt.Errorf("no enigma machine")
	}

	return what.machine.DecryptWithPlugBoard(cipherText, key, plugBoard)
}

func (what *Enigma) DecryptWithSetting(cipherText []byte, setting *settings.Setting) ([]byte, error) {
	if what.machine == nil {
		return nil, fmt.Errorf("no enigma machine")
	}

	return what.machine.DecryptWithSetting(cipherText, setting)
}

func (what *Enigma) GenerateKey() (string, error) {
	if what.machine == nil {
		return "", fmt.Errorf("no enigma machine")
	}

	return what.machine.GenerateKey()
}

func (what *Enigma) Sanitize(plainText string) (string, error) {
	if what.machine == nil {
		return "", fmt.Errorf("no enigma machine")
	}

	return what.machine.Sanitize(plainText), nil
}
