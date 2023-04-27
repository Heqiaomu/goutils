package util

import (
	"crypto/md5"
	"encoding/hex"
)

//MD5 get 32 bit md5 hash
func MD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

//MD5Part return part([start : end]) of md5 hash
func MD5Part(text string, start, end int) string {
	hash := md5.Sum([]byte(text))
	tr := hex.EncodeToString(hash[:])
	var rHash string
	if start > end || start > 32 || end > 32 {
		rHash = tr[:]
	} else {
		rHash = tr[start:end]
	}

	return rHash
}
