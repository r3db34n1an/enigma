package enigma

import (
	_ "embed"
	"fmt"
	"github.com/r3db34n1an/enigma/pkg/defs"
	"github.com/r3db34n1an/enigma/pkg/settings"
	"regexp"
	"strings"
	"unicode"
)

// Sources:
//	- https://blog.gopheracademy.com/advent-2016/enigma-emulator-in-go/
//	- https://www.ciphermachinesandcryptology.com/en/enigmaproc.htm
//  - https://www.codesandciphers.org.uk/enigma/index.htm
//	- https://en.wikipedia.org/wiki/Enigma_machine
//	- https://www.cryptomuseum.com/crypto/enigma/index.htm

type Enigma struct {
	preserveFormatting bool
	preserveCase       bool
	setting            settings.Setting
}

func NewEnigma(preserveFormatting bool, preserveCase bool) (*Enigma, error) {
	return &Enigma{
		preserveFormatting: preserveFormatting,
		preserveCase:       preserveCase,
	}, nil
}

func (what *Enigma) Init() error {
	return nil
}

func (what *Enigma) Encrypt(plainText []byte, key string) ([]byte, error) {
	return what.EncryptWithPlugBoard(plainText, key, "")
}

func (what *Enigma) EncryptWithSetting(plainText []byte, setting *settings.Setting) ([]byte, error) {
	*what = Enigma{
		preserveFormatting: what.preserveFormatting,
		preserveCase:       what.preserveCase,
	}

	if setting == nil {
		return nil, fmt.Errorf("no setting")
	}

	what.setting = *setting
	var cipherText []byte
	for index, plain := range plainText {
		if what.preserveFormatting {
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

		plainRune := rune(plain)
		encrypted := strings.IndexRune(defs.UpperCase, unicode.ToUpper(plainRune))
		if encrypted < 0 {
			return nil, fmt.Errorf("invalid character %q", plain)
		}

		encrypted = what.setting.PlugBoard.Transform(encrypted)
		if encrypted < 0 || encrypted > len(defs.UpperCase) {
			return nil, fmt.Errorf("plug board encryption of %q failed", plain)
		}

		encrypted = what.setting.Rotors.Encrypt(encrypted)
		if encrypted < 0 || encrypted > len(defs.UpperCase) {
			return nil, fmt.Errorf("plug rotor encryption of %q failed", plain)
		}

		encrypted = what.setting.Reflector.Reflect(encrypted)
		if encrypted < 0 || encrypted > len(defs.UpperCase) {
			return nil, fmt.Errorf("plug reflection of %q failed", plain)
		}

		encrypted = what.setting.Rotors.Decrypt(encrypted)
		if encrypted < 0 || encrypted > len(defs.UpperCase) {
			return nil, fmt.Errorf("plug rotor decryption of %q failed", plain)
		}

		encrypted = what.setting.PlugBoard.Transform(encrypted)
		if encrypted < 0 || encrypted >= len(defs.UpperCase) {
			return nil, fmt.Errorf("plug board decryption of %q failed", plain)
		}

		encryptedRune := rune(defs.UpperCase[encrypted])
		if unicode.IsLower(plainRune) && what.preserveCase {
			encryptedRune = unicode.ToLower(encryptedRune)
		}

		cipherText = append(cipherText, byte(encryptedRune))
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

func (what *Enigma) DecryptWithSetting(cipherText []byte, setting *settings.Setting) ([]byte, error) {
	*what = Enigma{
		preserveFormatting: what.preserveFormatting,
		preserveCase:       what.preserveCase,
	}

	if setting == nil {
		return nil, fmt.Errorf("no setting")
	}

	what.setting = *setting
	var plainText []byte
	for index, encrypted := range cipherText {
		if what.preserveFormatting {
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

		encryptedRune := rune(encrypted)
		plain := strings.IndexRune(defs.UpperCase, unicode.ToUpper(encryptedRune))
		if plain < 0 {
			return nil, fmt.Errorf("invalid character %q", encrypted)
		}

		plain = what.setting.PlugBoard.Transform(plain)
		if plain < 0 || plain > len(defs.UpperCase) {
			return nil, fmt.Errorf("plug board encryption of %q failed", encrypted)
		}

		plain = what.setting.Rotors.Encrypt(plain)
		if plain < 0 || plain > len(defs.UpperCase) {
			return nil, fmt.Errorf("plug rotor encryption of %q failed", encrypted)
		}

		plain = what.setting.Reflector.Reflect(plain)
		if plain < 0 || plain > len(defs.UpperCase) {
			return nil, fmt.Errorf("plug reflection of %q failed", encrypted)
		}

		plain = what.setting.Rotors.Decrypt(plain)
		if plain < 0 || plain > len(defs.UpperCase) {
			return nil, fmt.Errorf("plug rotor decryption of %q failed", encrypted)
		}

		plain = what.setting.PlugBoard.Transform(plain)
		if plain < 0 || plain >= len(defs.UpperCase) {
			return nil, fmt.Errorf("plug board decryption of %q failed", encrypted)
		}

		plainRune := rune(defs.UpperCase[plain])
		if unicode.IsLower(encryptedRune) && what.preserveCase {
			plainRune = unicode.ToLower(plainRune)
		}

		plainText = append(plainText, byte(plainRune))
	}

	return plainText, nil
}

func (what *Enigma) DecryptWithPlugBoard(cipherText []byte, key string, plugBoard string) ([]byte, error) {
	setting, keyError := what.readKeyAndPlugBoard(key, plugBoard)
	if keyError != nil {
		return nil, fmt.Errorf("failed to read key: %v", keyError)
	}

	return what.DecryptWithSetting(cipherText, setting)
}

func (what *Enigma) GenerateKey() (string, error) {
	var setting settings.Setting
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

func (what *Enigma) readKeyAndPlugBoard(key string, plugBoard string) (*settings.Setting, error) {
	setting := new(settings.Setting)

	var importedKey settings.ExportSetting
	parseError := importedKey.Parse(key)
	if parseError != nil {
		return nil, fmt.Errorf("failed to parse key: %v", parseError)
	}

	importError := setting.Import(importedKey)
	if importError != nil {
		return nil, fmt.Errorf("failed to import key: %v", importError)
	}

	if len(plugBoard) > 0 {
		plugBoardError := setting.LoadPlugBoard(plugBoard)
		if plugBoardError != nil {
			return nil, fmt.Errorf("failed to load plug board: %v", plugBoardError)
		}
	}

	return setting, nil
}

func (what *Enigma) shouldEncrypt(c byte) bool {
	if what.preserveCase {
		return strings.ContainsRune(defs.UpperCase, unicode.ToUpper(rune(c)))
	}

	return strings.ContainsRune(defs.UpperCase, rune(c))
}
