package license

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"time"
)

const (
	licensePath = "./LICENSE"
	// emptyIP empty ip default placeholder
	EmptyIP = "emptyIPV4"

	OverTime = 60e9 //60s
	//query
	QueryPriv = "fd26a860237b461d1baec33274c0705f256845ea846a7f40f48175fbf2c52a95"
	//response
	ResponsePub = "35aaf83087eb67f795cc52eaf9b2a4a1fdfb72840a3fe38a4d84eca718c072b4e0afc599f099c2381d26edbc362bf9bfe7eb6253c6ed795db6e5590dd31033a6"
)

const (
	//VerificationCycle Verification Cycle
	VerificationCycle         = time.Hour
	LimitOfRequestsInOneCycle = 5
)

// ConvertPublicKey public key constant, this public key will used to verify the license
func ConvertPublicKey(pubKey string) (key *ecdsa.PublicKey, err error) {

	if len(pubKey) != 128 {
		errors.New("invalid public key")
		return nil, err
	}

	a := pubKey[:64]
	b := pubKey[64:]

	key = new(ecdsa.PublicKey)
	key.X, _ = big.NewInt(0).SetString(a, 16)
	key.Y, _ = big.NewInt(0).SetString(b, 16)

	return key, nil
}
