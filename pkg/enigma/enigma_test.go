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
	{
		Key: `
rotors:
     - name: II
       position: Z
       ring_setting: F
     - name: IV
       position: D
       ring_setting: Q
     - name: I
       position: "N"
       ring_setting: Z
reflector: B
plug_board:
     A: C
     B: Q
     F: J
     K: O
     L: S
     M: "Y"
     "N": W
     P: Z
     R: T
     U: V
`,
		Plain: `TO MY TRUSTED ALLY:

THE MACHINES HERE ARE BRILLIANT AND RELENTLESS. THE HACKERS HAVE USED THEM TO ENCODE THEIR PLANS
THINKING THEY ARE UNBREAKABLE, BUT THEY UNDERESTIMATE US. HIDDEN WITHIN THEIR TANGLED WIRES AND
GEARS IS THE KEY TO THE NEXT STEP. THE CODE YOU SEEK IS LOCATED HERE:

HTTPS://WWW.GUTENBERG.ORG/CACHE/EPUB/35/PG35.TXT WITH OFFSET TWO HUNDRED SEVENTY-EIGHT

MBGNJIGPNJMEMKKLLJPAOHIFMMKHEPIJ

YOUR DIRECTOR OF APPSEC
`,
		Encrypted: `FB IG HTNRZAQ KKOA:

SRG LCBNQWUW UFSM LSX EJTHUPWFJ PTS FVDDSVMSFT. YDF PBUEMWW TCAZ YYYS QZJB DY HEPAKK LPCEA OTJQB
MDFLCMJX FSZT XMJ QJPEXGJCAOR, FHM LQMK SEHRYWOLXWNAN HU. CQFJTW CMGAMB GCVZC XRQIQZM GLVJB IXM
IWBTQ LQ UTH LUJ AF RVJ VICJ KFQW. WRC LNEN JUI CFAE KO ADYLNDU KGKH:

QGOJW://RLP.FPDLVCDSD.ZFK/YDTOM/JYVJ/35/XO35.XHV LVND YLOYRY AZN PBDWOBU ODZZBBW-YDASA

EACWITXBDTZWADPRCWGNDOTSHBVJUVJK

JEYJ KHERHYBK BZ VQJQVR
`,
		Decrypted: `TO MY TRUSTED ALLY:

THE MACHINES HERE ARE BRILLIANT AND RELENTLESS. THE HACKERS HAVE USED THEM TO ENCODE THEIR PLANS
THINKING THEY ARE UNBREAKABLE, BUT THEY UNDERESTIMATE US. HIDDEN WITHIN THEIR TANGLED WIRES AND
GEARS IS THE KEY TO THE NEXT STEP. THE CODE YOU SEEK IS LOCATED HERE:

HTTPS://WWW.GUTENBERG.ORG/CACHE/EPUB/35/PG35.TXT WITH OFFSET TWO HUNDRED SEVENTY-EIGHT

MBGNJIGPNJMEMKKLLJPAOHIFMMKHEPIJ

YOUR DIRECTOR OF APPSEC
`,
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
