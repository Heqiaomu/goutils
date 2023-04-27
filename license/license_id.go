package license

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrCodeVerifyFail verify failed error
	ErrCodeVerifyFail = errors.New("code verify fail")
	// ErrCodeIsNotBase64 error code is not base64
	ErrCodeIsNotBase64 = errors.New("code is not base 64")
)

// LicenseCodeContentInner Inner license code definition
type LicenseCodeContentInner struct {
	UID       string
	EndTime   time.Time //UNIX time
	Version   string    //commercial version
	Online    bool
	LicenseID int64  `asn1:"omitempty,optional"`
	Extra     string `asn1:"omitempty,optional"`
}

// LicenseCodeContent outer license code content definition
type LicenseCodeContent struct {
	UID          *UID
	EndTime      time.Time //UNIX time
	Version      string    //commercial version
	Online       bool
	LicenseID    int64  `asn1:"omitempty,optional"`
	Extra        string `asn1:"omitempty,optional"`
	VerifyDomain string `asn1:"omitempty,optional"` //只有Online为true时起作用
}

// DecodeLicenseCode process
func DecodeLicenseCode(b string, key *ecdsa.PublicKey) (*LicenseCodeContent, error) {
	code, err := base64.URLEncoding.DecodeString(b)
	if err != nil {
		return nil, ErrCodeIsNotBase64
	}

	sign := code[:66]
	r := big.NewInt(0).SetBytes(sign[:32])
	s := big.NewInt(0).SetBytes(sign[32:64])

	sha := sha256.New()
	sha.Write(code[66:])
	key.Curve = elliptic.P256()
	if !ecdsa.Verify(key, sha.Sum(nil), r, s) {
		return nil, ErrCodeVerifyFail
	}

	h := code[66:]
	result := new(LicenseCodeContentInner)

	_, err = asn1.Unmarshal(h, result)
	if err != nil {
		return nil, err
	}
	tmp, err := strconv.ParseUint(result.UID[:2], 16, 8)
	if err != nil {
		return nil, err
	}
	uid := Decode(result.UID)
	if uid == nil || UIDType(tmp) != uid.T {
		return nil, errors.New("err")
	}
	return &LicenseCodeContent{
		UID:       uid,
		EndTime:   result.EndTime,
		Version:   result.Version,
		Online:    result.Online,
		LicenseID: result.LicenseID,
		Extra:     result.Extra,
	}, nil
}

func (lc *LicenseCodeContent) CheckLicense() (time.Time, error) {
	ex := strings.Split(lc.Extra, "|")
	client := ""
	prod := ""
	if len(ex) < 2 {
		client = "hyperchain"
		prod = "趣链科技基础平台部"
	} else {
		client = strings.TrimSpace(ex[1])
		prod = "趣链科技基础平台部"
	}

	if time.Now().After(lc.EndTime) {
		fmt.Println("[license] license is expired")
		return time.Unix(0, 0), ErrExpired
	}

	if !lc.UID.Verify() {
		fmt.Println(`[license] license uid(maybe ip address) error`)
		return time.Unix(0, 0), ErrUID
	}

	fmt.Printf("[license] %v: license to %v, exp data %v\n",
		time.Now().Format("15:04:05"), client, lc.EndTime.Format("2006-01-02"))

	storeLC = lc
	storeClient = client
	storeProd = prod
	return lc.EndTime, nil
}
