package crypt

import (
	"net"
)

// cxncode 连接码(connection code)
// 连接码本质上是一个由字母集[0-9A-Za-z]构成的62进制数
// 字母表 48-57对应0-9；65-90对应A-Z；97-122对应a-z

// IPToCxnCode ip转连接码
// 将IP格式的IP转化成6位连接码
func IPToCxnCode(ip string) string {
	// 校验IP的格式是否正确，支持IPv4 IPv6
	nIP := net.ParseIP(ip)
	if nIP == nil {
		return ""
	}
	// 将IP转成IPv4格式
	nIP = nIP.To4()
	if nIP == nil {
		return ""
	}
	// 将IP转化为10进制数字aIP
	aIP := uint32(nIP[0])<<24 + uint32(nIP[1])<<16 + uint32(nIP[2])<<8 + uint32(nIP[3])
	// 将aIP转化成6位连接码
	cxn := make([]byte, 6)
	for i := 0; i < 6; i++ {
		cxn[i] = 48
	}
	var rem uint32
	cxni := 5
	for aIP != 0 {
		rem = aIP % 62
		if rem <= 9 {
			cxn[cxni] = byte(48 + rem)
		} else if rem <= 34 {
			cxn[cxni] = byte(65 + rem - 10)
		} else {
			cxn[cxni] = byte(97 + rem - 36)
		}
		cxni--
		aIP = aIP / 62
	}
	return string(cxn)
}

// CxnCodeToIPv4 连接码转IPv4格式的字符串
func CxnCodeToIPv4(cxn string) string {
	if len(cxn) != 6 {
		return ""
	}

	var aIP uint32
	for i, c := range cxn {
		var n uint32
		var p uint32 = 1
		if c >= 48 && c <= 57 {
			n = uint32(c - 48)
		} else if c >= 65 && c <= 90 {
			n = uint32(c-65) + 10
		} else if c >= 97 && c <= 122 {
			n = uint32(c-97) + 36
		} else {
			return ""
		}
		for j := 0; j < 5-i; j++ {
			p *= 62
		}
		aIP += n * p
	}

	return net.IPv4(byte((aIP>>24)&0xFF), byte((aIP>>16)&0xFF), byte((aIP>>8)&0xFF), byte(aIP&0xFF)).String()
}
