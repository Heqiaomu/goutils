package license

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var (
	ErrExpired     = errors.New("license expired")
	ErrUID         = errors.New("uid error")
	ErrCheckOnline = errors.New("check online failed")
	ErrSignature   = errors.New("signature err")
	ErrReadFail    = errors.New("read license error")
	ErrSyntax      = errors.New("license syntax error")
	ErrUnknown     = errors.New("unknown license err, maybe panic")
	ErrExtern      = errors.New("extern check error")
	ErrExternPanic = errors.New("extern check panic")
)

const MaxLimit = 24 * 5

type ExtentVerifyFunc func(*LicenseCodeContent) error

func CheckLicense(exit chan bool, pubKey string, extent ExtentVerifyFunc) {

	key, err := ConvertPublicKey(pubKey)
	if err != nil {
		log.Printf("License Verify Failed: %v", err)
		return
	}

	//check at start up
	ierr, _ := IsLicenseExpired(key, extent)
	if ierr != nil {
		log.Printf("License Verify Failed At Startup: %v", ierr)
		notifySystemExit(exit)
		return
	}
	c <- true
	// this ensures that license checker always hit in `os thread` to avoid jmuping to other threads
	// since in this approach, working directory will not be affected by other operators.
	runtime.LockOSThread()
	// check license immediately once hyperchain start.
	timer := time.NewTimer(VerificationCycle)
	i := MaxLimit
	for {
		select {
		case <-timer.C:
			ierr, _ := IsLicenseExpired(key, extent)
			if ierr != nil {
				i--
				if i > 0 {
					continue
				}
				log.Printf("License Verify Failed: %v", ierr)
				notifySystemExit(exit)
				return
			}
			log.Printf("License Verify Pass!")
			timer.Reset(VerificationCycle)
		}
	}
}

// ReadLicense read the license byte
func ReadLicense() (license []byte, err error) {
	f, err := os.Open(licensePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("license file not found: %s ", licensePath)
		}
		return nil, fmt.Errorf("read license error: %s", err.Error())
	}
	return ioutil.ReadAll(f)
}

// IsLicenseExpired - check whether license is expired.
func IsLicenseExpired(key *ecdsa.PublicKey, extent ExtentVerifyFunc) (err error, msg string) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrUnknown
			msg = fmt.Sprintf("Invalid License:%v", r)
			log.Printf(msg)
		}
	}()

	//1.read license
	license, rerr := ReadLicense()
	if rerr != nil {
		log.Printf("reading license failed, %s", rerr.Error())
		err = ErrReadFail
		msg = rerr.Error()
		return
	}
	pattern, _ := regexp.Compile("Identification: (.*)")
	identification := pattern.FindString(string(license))[16:]
	identification = strings.TrimSpace(identification)
	//2.parse license
	licenseID, derr := DecodeLicenseCode(identification, key)
	if derr != nil {
		err = ErrSyntax
		if derr == ErrCodeVerifyFail {
			err = ErrSignature
		}
		msg = derr.Error()
		return
	}

	//3. extern verify
	func() {
		defer func() {
			r := recover()
			if !reflect.DeepEqual(r, nil) {
				err = ErrExternPanic
				msg = fmt.Sprint(r)
			}
		}()
		if extent != nil && extent(licenseID) != nil {
			log.Printf("extern check error:" + err.Error())
			err = ErrExtern
			msg = err.Error()
			return
		}
	}()

	//4. check license
	expiredTime, cerr := licenseID.CheckLicense()
	if cerr != nil {
		log.Printf(cerr.Error())
		err = cerr
		msg = cerr.Error()
		return
	}
	//5. check license online
	if licenseID.Online {
		log.Printf("check online :" + licenseID.VerifyDomain)
		qerr := query(identification, licenseID.VerifyDomain, key)
		if qerr != nil {
			err = ErrCheckOnline
			msg = qerr.Error()
			return
		}
		log.Printf("check online success")
	}

	err = nil
	msg = expiredTime.String()
	return
}

// RetrieveLicenseUID retrieve license binding UID
func RetrieveLicenseUID() (ID string, err error, key *ecdsa.PublicKey) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("invalid license: %s", r)
		}
	}()

	license, err := ReadLicense()
	if err != nil {
		err = fmt.Errorf("invalid license: %s", err)
		return
	}

	pattern, _ := regexp.Compile("Identification: (.*)")
	identification := pattern.FindString(string(license))[16:]
	identification = strings.TrimSpace(identification)

	lc, err := DecodeLicenseCode(identification, key)
	if err != nil {
		err = fmt.Errorf("invalid license: %s", err)
		return
	}
	ID = lc.UID.ToString()
	return
}

// notifySystemExit - license expired or not found, shut down system.
func notifySystemExit(exit chan bool) {
	exit <- true

	<-time.After(3 * time.Second)
	//common.ExitWithLockFile(false, "License expired!")
}
