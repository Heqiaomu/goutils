package license

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

var queryPriv *ecdsa.PrivateKey
var responsePub *ecdsa.PublicKey

func init() {
	queryPriv = new(ecdsa.PrivateKey)
	queryPriv.D, _ = big.NewInt(0).SetString(QueryPriv, 16)
	queryPriv.Curve = elliptic.P256()

	responsePub = new(ecdsa.PublicKey)
	responsePub.X, _ = big.NewInt(0).SetString(ResponsePub[:64], 16)
	responsePub.Y, _ = big.NewInt(0).SetString(ResponsePub[64:], 16)
	responsePub.Curve = elliptic.P256()
}

type request struct {
	LicenseID int64  `json:"id"`
	Time      int64  `json:"time"`
	Nonce     string `json:"nonce"` //S1(license[20]+time),base64
}

type response struct {
	Time     int64  `json:"time"`
	Response string `json:"response"` //S2(license[20]+license[20]+time),base64
}

//return err if get a non-200 response
func query(i, s string, key *ecdsa.PublicKey) error {
	//解析license
	l, derr := DecodeLicenseCode(i, key)
	if derr != nil {
		return errors.New("pleas use license after version 1.7")
	}
	//请求
	startTime := time.Now()
	startTimeUnix := startTime.Unix()
	hash := i[:20] + strconv.FormatInt(startTimeUnix, 10)
	cdfda := sha256Hash(hash)
	dsa := ecdsaSign(queryPriv, cdfda)
	rde, _ := json.Marshal(request{
		LicenseID: l.LicenseID,
		Time:      startTimeUnix,
		Nonce:     dsa,
	})

	rr, serr := httpsPost(s, bytes.NewReader(rde))
	if serr != nil {
		return serr
	}

	rep := new(response)
	uerr := json.Unmarshal(rr, rep)
	if uerr != nil {
		return uerr
	}

	hash = string(i)[:20] + hash
	cdfda = sha256Hash(hash)
	b, err := ecdsaVerify(responsePub, cdfda, rep.Response)
	if !b || err != nil {
		return errors.New("verify response fail")
	}

	endTime := time.Unix(rep.Time, 0)
	t := endTime.Sub(startTime).Nanoseconds()
	if t > OverTime || t < -OverTime {
		return errors.New("response time out")
	}
	return nil
}

// ECDSASignature represents an ECDSA signature
type ecdsaSignature struct {
	R, S *big.Int
}

func ecdsaVerify(verKey *ecdsa.PublicKey, msg []byte, signature string) (bool, error) {
	e := errors.New("签名不正确")
	s, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		return false, e
	}
	ecdsaSignature := new(ecdsaSignature)
	_, err = asn1.Unmarshal(s, ecdsaSignature)
	if err != nil {
		return false, e
	}
	return ecdsa.Verify(verKey, msg, ecdsaSignature.R, ecdsaSignature.S), nil
}

func ecdsaSign(verKey *ecdsa.PrivateKey, msg []byte) string {
	r, err := verKey.Sign(rand.Reader, msg, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	return base64.URLEncoding.EncodeToString(r)
}

//Hash 计算hash
func sha256Hash(s string) []byte {
	in := []byte(s)
	sha := sha256.New()
	sha.Write(in)
	return sha.Sum(nil)
}

func httpsPost(url string, body io.Reader) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, gerr := client.Post(url, "application/json", body)
	if gerr != nil {
		return nil, errors.New("validate err")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("response state is not 200")
	}
	return ioutil.ReadAll(resp.Body)
}
