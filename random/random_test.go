package random

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRandom(t *testing.T) {
	Convey("random test", t, func() {
		a := GetRandomString(6)
		b := GetRandomString(6)
		So(a, ShouldNotEqual, b)
	})
}
