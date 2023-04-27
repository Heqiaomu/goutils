package crypt

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCxnCode(t *testing.T) {
	Convey(`test cxn code`, t, func() {
		cxnCode := IPToCxnCode("10.1.41.111")
		ShouldEqual(cxnCode, "0BMH0h")
		ip := CxnCodeToIPv4(cxnCode)
		ShouldEqual(ip, "10.1.41.111")

		cxnCode = IPToCxnCode("0.0.0.0")
		ShouldEqual(cxnCode, "000000")
		ip = CxnCodeToIPv4(cxnCode)
		ShouldEqual(ip, "0.0.0.0")

		cxnCode = IPToCxnCode("255.255.255.255")
		ShouldEqual(cxnCode, "4gfFC3")
		ip = CxnCodeToIPv4(cxnCode)
		ShouldEqual(ip, "255.255.255.255")

		cxnCode = IPToCxnCode("192.200.17.1")
		ShouldEqual(cxnCode, "3WswoD")
		ip = CxnCodeToIPv4(cxnCode)
		ShouldEqual(ip, "192.200.17.1")

		cxnCode = IPToCxnCode("&+-//")
		ShouldEqual(cxnCode, "")

		ip = CxnCodeToIPv4("&+-//")
		ShouldEqual(ip, "")
	})
}
