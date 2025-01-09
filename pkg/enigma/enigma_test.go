package enigma

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCases = []TestCase{
	{
		Key: `
rotors:
    - name: V
      position: Z
      ring_setting: L
    - name: I
      position: H
      ring_setting: "Y"
    - name: II
      position: J
      ring_setting: H
reflector: B
plug_board:
    D: Z
    G: B
    M: E
    "N": L
    P: "Y"
    R: Q
    S: J
    T: F
    W: I
    X: V
`,
		Plain:     "OURDIRECTORISSAFE",
		Encrypted: "HRMZLDHBJRFRJXMAH",
		Decrypted: "OURDIRECTORISSAFE",
	},
	{
		Key: `
rotors:
    - name: III
      position: A
      ring_setting: A
    - name: II
      position: B
      ring_setting: A
    - name: IV
      position: C
      ring_setting: A
reflector: B
plug_board:
    A: B
    C: D
    E: F
`,
		Plain:     "HELLOWORLD",
		Encrypted: "YGMGTTPJNJ",
		Decrypted: "HELLOWORLD",
	},
	{
		Key: `
rotors:
    - name: III
      position: A
      ring_setting: E
    - name: II
      position: A
      ring_setting: J
    - name: IV
      position: A
      ring_setting: R
reflector: B
plug_board:
    A: E
    D: Q
    R: C
    V: B
    M: T
    O: G
    P: F
    Y: L
    J: W
    I: Z
`,
		Plain:     "NEVERXGONNAXGIVEXYOUXUPXNEVERXGONNAXLETXYOUXDOWNXIXBETXYOUXHAVEXNEVERXBEENXRICKROLLEDXWITHXANXENIGMAXBEFOREXGOODXLUCKXANDXTHANKXYOUXFORXREADING",
		Encrypted: "PJHLFULLUECCPFLCIVPMFDAWJCWANLVXAIXFHMACNLVNCSXOIXFUTGWXSRULRTXPOIPUINCYOGWKGZAZDMVPOUIDCRSCHSZCNTFJADAVIKOGSYAJGAFNELPOMBMTXEXVAREVMSBNHLJFEGZ",
		Decrypted: "NEVERXGONNAXGIVEXYOUXUPXNEVERXGONNAXLETXYOUXDOWNXIXBETXYOUXHAVEXNEVERXBEENXRICKROLLEDXWITHXANXENIGMAXBEFOREXGOODXLUCKXANDXTHANKXYOUXFORXREADING",
	},
}

func TestGenerate(t *testing.T) {
	for _, item := range testCases {
		for i := 0; i < 1000; i++ {
			cipher, cipherError := NewEnigma(true)
			assert.Nil(t, cipherError)
			if cipher == nil {
				return
			}

			assert.Nil(t, cipher.Init())

			key, keyError := cipher.GenerateKey()
			assert.Nil(t, keyError)
			assert.NotEmpty(t, key)

			encrypted, encryptError := cipher.Encrypt([]byte(item.Plain), item.Key)
			assert.Nil(t, encryptError)

			decrypted, decryptError := cipher.Decrypt(encrypted, item.Key)
			assert.Nil(t, decryptError)
			assert.Equal(t, item.Plain, string(decrypted))
		}
	}
}

func TestEncrypt(t *testing.T) {
	for _, item := range testCases {
		cipher, cipherError := NewEnigma(true)
		assert.Nil(t, cipherError)
		if cipher == nil {
			return
		}

		assert.Nil(t, cipher.Init())

		result, encryptError := cipher.Encrypt([]byte(item.Plain), item.Key)
		assert.Nil(t, encryptError)
		assert.Equal(t, item.Encrypted, string(result))
	}
}

func TestDecrypt(t *testing.T) {
	for _, item := range testCases {
		cipher, cipherError := NewEnigma(true)
		assert.Nil(t, cipherError)
		if cipher == nil {
			return
		}

		assert.Nil(t, cipher.Init())

		result, decryptError := cipher.Decrypt([]byte(item.Encrypted), item.Key)
		assert.Nil(t, decryptError)
		assert.Equal(t, item.Decrypted, string(result))
	}
}
