package defs

import (
	crand "crypto/rand"
	"math/big"
	mrand "math/rand"
)

func RandomInt(min int, max int) int {
	if min > max {
		min, max = max, min
	}

	delta := max - min + 1
	bigNumber, randError := crand.Int(crand.Reader, big.NewInt(int64(delta)))
	if randError != nil {
		return mrand.Int()%delta + min // #nosec:G404
	}

	bigNumber.Add(bigNumber, big.NewInt(int64(min)))
	return int(bigNumber.Int64())
}
