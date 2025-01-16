package enigma

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"
)

// Sources:
//	- https://blog.gopheracademy.com/advent-2016/enigma-emulator-in-go/
//	- https://www.ciphermachinesandcryptology.com/en/enigmaproc.htm
//  - https://www.codesandciphers.org.uk/enigma/index.htm
//	- https://en.wikipedia.org/wiki/Enigma_machine
//	- https://www.cryptomuseum.com/crypto/enigma/index.htm

//go:embed config/settings.yaml
var settingsYaml []byte

//go:embed config/rotors.yaml
var rotorsYaml []byte

//go:embed config/reflectors.yaml
var reflectorsYaml []byte

const (
	upperCase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type Enigma struct {
	copyExtra bool
	setting   Setting
}

func NewEnigma(copyExtra bool) (*Enigma, error) {
	return &Enigma{
		copyExtra: copyExtra,
	}, nil
}

func (what *Enigma) Init() error {
	return nil
}

func (what *Enigma) Encrypt(plainText []byte, key string) ([]byte, error) {
	return what.EncryptWithPlugBoard(plainText, key, "")
}

func (what *Enigma) EncryptWithSetting(plainText []byte, setting *Setting) ([]byte, error) {
	*what = Enigma{
		copyExtra: what.copyExtra,
	}

	if setting == nil {
		return nil, fmt.Errorf("no setting")
	}

	what.setting = *setting
	var cipherText []byte
	for index, plain := range plainText {
		if what.copyExtra {
			if !what.shouldEncrypt(plain) {
				cipherText = append(cipherText, plain)
				continue
			}
		} else {
			if index > 0 {
				if index%5 == 0 {
					cipherText = append(cipherText, ' ')
				}

				if index%80 == 0 {
					cipherText = append(cipherText, '\n')
				}
			}
		}

		what.setting.Rotors.Move()

		encrypted := strings.IndexRune(upperCase, rune(plain))
		if encrypted < 0 {
			return nil, fmt.Errorf("invalid character %q", plain)
		}

		encrypted = what.setting.PlugBoard.Transform(encrypted)
		if encrypted < 0 || encrypted > len(upperCase) {
			return nil, fmt.Errorf("plug board encryption of %q failed", plain)
		}

		encrypted = what.setting.Rotors.Encrypt(encrypted)
		if encrypted < 0 || encrypted > len(upperCase) {
			return nil, fmt.Errorf("plug rotor encryption of %q failed", plain)
		}

		encrypted = what.setting.Reflector.Reflect(encrypted)
		if encrypted < 0 || encrypted > len(upperCase) {
			return nil, fmt.Errorf("plug reflection of %q failed", plain)
		}

		encrypted = what.setting.Rotors.Decrypt(encrypted)
		if encrypted < 0 || encrypted > len(upperCase) {
			return nil, fmt.Errorf("plug rotor decryption of %q failed", plain)
		}

		encrypted = what.setting.PlugBoard.Transform(encrypted)
		if encrypted < 0 || encrypted >= len(upperCase) {
			return nil, fmt.Errorf("plug board decryption of %q failed", plain)
		}

		cipherText = append(cipherText, upperCase[encrypted])
	}

	return cipherText, nil
}

func (what *Enigma) EncryptWithPlugBoard(plainText []byte, key string, plugBoard string) ([]byte, error) {
	setting, keyError := what.readKeyAndPlugBoard(key, plugBoard)
	if keyError != nil {
		return nil, fmt.Errorf("failed to read key: %v", keyError)
	}

	return what.EncryptWithSetting(plainText, setting)
}

func (what *Enigma) Decrypt(cipherText []byte, key string) ([]byte, error) {
	return what.DecryptWithPlugBoard(cipherText, key, "")
}

func (what *Enigma) DecryptWithSetting(cipherText []byte, setting *Setting) ([]byte, error) {
	*what = Enigma{
		copyExtra: what.copyExtra,
	}

	if setting == nil {
		return nil, fmt.Errorf("no setting")
	}

	what.setting = *setting
	var plainText []byte
	for index, encrypted := range cipherText {
		if what.copyExtra {
			if !what.shouldEncrypt(encrypted) {
				plainText = append(plainText, encrypted)
				continue
			}
		} else {
			if index > 0 {
				if index%5 == 0 {
					plainText = append(plainText, ' ')
				}

				if index%80 == 0 {
					plainText = append(plainText, '\n')
				}
			}
		}

		what.setting.Rotors.Move()

		plain := strings.IndexRune(upperCase, rune(encrypted))
		if plain < 0 {
			return nil, fmt.Errorf("invalid character %q", encrypted)
		}

		plain = what.setting.PlugBoard.Transform(plain)
		if plain < 0 || plain > len(upperCase) {
			return nil, fmt.Errorf("plug board encryption of %q failed", encrypted)
		}

		plain = what.setting.Rotors.Encrypt(plain)
		if plain < 0 || plain > len(upperCase) {
			return nil, fmt.Errorf("plug rotor encryption of %q failed", encrypted)
		}

		plain = what.setting.Reflector.Reflect(plain)
		if plain < 0 || plain > len(upperCase) {
			return nil, fmt.Errorf("plug reflection of %q failed", encrypted)
		}

		plain = what.setting.Rotors.Decrypt(plain)
		if plain < 0 || plain > len(upperCase) {
			return nil, fmt.Errorf("plug rotor decryption of %q failed", encrypted)
		}

		plain = what.setting.PlugBoard.Transform(plain)
		if plain < 0 || plain >= len(upperCase) {
			return nil, fmt.Errorf("plug board decryption of %q failed", encrypted)
		}

		plainText = append(plainText, upperCase[plain])
	}

	return plainText, nil
}

func (what *Enigma) DecryptWithPlugBoard(cipherText []byte, key string, plugBoard string) ([]byte, error) {
	setting, keyError := what.readKeyAndPlugBoard(key, plugBoard)
	if keyError != nil {
		return nil, fmt.Errorf("failed to read key: %v", keyError)
	}

	return what.EncryptWithSetting(cipherText, setting)
}

func (what *Enigma) GenerateKey() (string, error) {
	var setting Setting
	randomError := setting.Random()
	if randomError != nil {
		return "", fmt.Errorf("failed to generate random setting: %v", randomError)
	}

	exported := setting.Export()
	generateError := exported.Generate()
	if generateError != nil {
		return "", fmt.Errorf("failed to generate key info: %v", generateError)
	}

	return exported.Key, nil
}

func (what *Enigma) Sanitize(plainText string) string {
	// historical letter substitutions in actual messages
	// `?` => `L`
	// `,` => `Y`
	// `.` => `X`
	// `:` => `XX`
	// `"` => `J`
	// `(`, ')' => `KK`

	plainText = strings.ReplaceAll(plainText, "?", "L")
	plainText = strings.ReplaceAll(plainText, ",", "Y")
	plainText = strings.ReplaceAll(plainText, ".", "X")
	plainText = strings.ReplaceAll(plainText, ":", "XX")
	plainText = strings.ReplaceAll(plainText, "\"", "J")
	plainText = strings.ReplaceAll(plainText, "(", "KK")
	plainText = strings.ReplaceAll(plainText, ")", "KK")

	re, regexError := regexp.Compile(`[^A-Z]`)
	if regexError != nil {
		return plainText
	}

	return re.ReplaceAllString(plainText, "")
}

func (what *Enigma) readKeyAndPlugBoard(key string, plugBoard string) (*Setting, error) {
	setting := new(Setting)

	var importedKey ExportSetting
	parseError := importedKey.Parse(key)
	if parseError != nil {
		return nil, fmt.Errorf("failed to parse key: %v", parseError)
	}

	importError := setting.Import(importedKey)
	if importError != nil {
		return nil, fmt.Errorf("failed to import key: %v", importError)
	}

	if len(plugBoard) > 0 {
		plugBoardError := setting.loadPlugBoard(plugBoard)
		if plugBoardError != nil {
			return nil, fmt.Errorf("failed to load plug board: %v", plugBoardError)
		}
	}

	return setting, nil
}

func (what *Enigma) shouldEncrypt(c byte) bool {
	return strings.ContainsRune(upperCase, rune(c))
}
