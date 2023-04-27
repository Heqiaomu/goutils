package license

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	"math/big"
	"time"
)

//challenge public key
var (
	x, _ = new(big.Int).SetString("ccce43de7f7e2c8c836f977fd3784d7a056acbb807bf0a502ab4dfeea2817161", 16)
	y, _ = new(big.Int).SetString("cc84dea9a25414a51a46f0c0d6b7af90cdd5647f3116ac32245da14009b1ded3", 16)
)

type BassChallenge struct {
	Nonce string `json:"nonce"` //base64
	Time  int64  `json:"time"`  //unix stamp
}

type BassResponse struct {
	Signature string `json:"signature"` //base64
}

func CheckExist(licenseServerURL string) error {
	randNum := make([]byte, 20)
	_, _ = rand.Read(randNum)
	query := &BassChallenge{
		Nonce: hex.EncodeToString(randNum),
		Time:  time.Now().Unix(),
	}
	queryBytes, _ := json.Marshal(*query)

	responseBytes, err := httpsPost(licenseServerURL, bytes.NewBuffer(queryBytes))
	if err != nil {
		return errors.Wrap(err, string(responseBytes))
	}
	response := new(BassResponse)
	err = json.Unmarshal(responseBytes, response)
	if err != nil {
		return err
	}

	hasher := sha256.New()
	hasher.Write(queryBytes)
	degist := hasher.Sum(nil)

	_, err = ecdsaVerify(&ecdsa.PublicKey{
		X:     x,
		Y:     y,
		Curve: elliptic.P256(),
	}, degist, response.Signature)
	if err != nil {
		return err
	}
	return nil
}
