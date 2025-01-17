package enigma

type TestCase struct {
	Key                string
	Plain              string
	Encrypted          string
	Decrypted          string
	PreserveCase       bool
	PreserveFormatting bool
}
