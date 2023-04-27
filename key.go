//Package util util
package util

import (
	"crypto/rsa"
	"io/ioutil"

	// "github.com/dgrijalva/jwt-go"

	"github.com/golang-jwt/jwt"
	"github.com/golang/glog"
)

// GeneratePublicKey generates public key
func GeneratePublicKey(filename string) *rsa.PublicKey {

	publicBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		glog.Errorf("Fail to generate public key, because no public key provide, err: [%v].", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		glog.Errorf("Fail to generate public key, because not a valid public key, err: [%v].", err)
	}
	return publicKey

}
