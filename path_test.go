package util

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/agiledragon/gomonkey"

	"bou.ke/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPathExistsMonkey(t *testing.T) {
	Convey("Test Path Exists util", t, func() {
		// patch os.Stat func
		guard := monkey.Patch(os.Stat, func(target string) (os.FileInfo, error) {
			switch target {
			case "exists.txt":
				return nil, nil // PathExists中未使用到fileinfo,直接返回nil,如适用到，根据分支情况进行mock
			case "notexists.txt":
				return nil, os.ErrNotExist
			case "unknowerr.txt":
				return nil, errors.New("unknowerr")
			}
			return nil, os.ErrNotExist
		})
		defer guard.Unpatch()

		Convey("exists test", func() {
			e, err := PathExists("exists.txt")
			So(e, ShouldBeTrue)
			So(err, ShouldBeNil)
		})
		Convey("not exists test", func() {
			ne, err := PathExists("notexists.txt")
			So(ne, ShouldBeFalse)
			So(err, ShouldBeNil)
		})
		Convey("unknow error test", func() {
			un, err := PathExists("unknowerr.txt")
			So(un, ShouldBeFalse)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestPathExistsGoMonkey(t *testing.T) {
	Convey("Test Path exists gomonkey", t, func() {
		pach := gomonkey.ApplyFunc(os.Stat, func(target string) (os.FileInfo, error) {
			switch target {
			case "exists.txt":
				return nil, nil // PathExists中未使用到fileinfo,直接返回nil,如适用到，根据分支情况进行mock
			case "notexists.txt":
				return nil, os.ErrNotExist
			case "unknowerr.txt":
				return nil, errors.New("unknowerr")
			}
			return nil, os.ErrNotExist
		})
		defer pach.Reset()

		Convey("exists test gomonkey", func() {
			e, err := PathExists("exists.txt")
			So(e, ShouldBeTrue)
			So(err, ShouldBeNil)
		})
		Convey("not exists test gomonkey", func() {
			ne, err := PathExists("notexists.txt")
			So(ne, ShouldBeFalse)
			So(err, ShouldBeNil)
		})
		Convey("unknow error test gomonkey", func() {
			un, err := PathExists("unknowerr.txt")
			So(un, ShouldBeFalse)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestPathExistsOri(t *testing.T) {
	Convey("Test PathExists Ori", t, func() {
		pwd, _ := os.Getwd()

		Convey("path exist", func() {
			e, err := PathExists(pwd)
			So(e, ShouldBeTrue)
			So(err, ShouldBeNil)
		})
		Convey("path not exist", func() {
			e, err := PathExists(filepath.Join(pwd, "notexists"))
			So(e, ShouldBeFalse)
			So(err, ShouldBeNil)
		})
	})
}
