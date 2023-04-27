package license

import (
	"errors"
	"net"
	"strconv"
)

type UIDType byte

var IDTypeName = [UIDTypeMaxLimit]string{
	TrailString,
	LocalIPString,
	GlobalIPString,
}

const (
	TrailString    = "trail"
	LocalIPString  = "localip"
	GlobalIPString = "globalip"
)

const (
	// InternalTrail type
	InternalTrail UIDType = iota
	// LocalIP type
	LocalIP
	// GlobalIP type
	GlobalIP
	//Ali
	//Tencent
	//Huawei
	//Microsoft
	//Amazon
	// UIDTypeMaxLimit type
	UIDTypeMaxLimit
)

var verifyUID = [UIDTypeMaxLimit]func(s string) bool{
	func(s string) bool {
		return s == EmptyIP || s == ""
	},
	checkIP,
	checkIP,
}

func checkIP(s string) bool {
	ipInLicense := net.ParseIP(s)
	if ipInLicense == nil {
		return false
	}
	ipv4, err := RetrieveLocalIPV4()
	if err != nil {
		return false
	}
	for i := range ipv4 {
		if net.ParseIP(ipv4[i]).Equal(ipInLicense) {
			return true
		}
	}
	return false
}

var CheckInput = [UIDTypeMaxLimit]func(s string) bool{
	//InternalTrail
	func(s string) bool {
		return s == EmptyIP || s == ""
	},
	//LocalIP
	func(s string) bool {
		ip := net.ParseIP(s)
		if !ip.IsGlobalUnicast() {
			return false
		}
		ip = ip.To4()
		if ip[0] == 10 || //A: 10.0.0.0 - 10.255.255.255
			(ip[0] == 172 && ip[1]&0xf0 == 0x10) || //B: 172.16.0.0 - 172.31.255.255
			(ip[0] == 192 && ip[1] == 168) { //C:192.168.0.0 - 192.168.255.255
			return true
		}
		return false
	},
	//GlobalIP
	func(s string) bool {
		ip := net.ParseIP(s)
		if ip == nil {
			return false
		}
		//127.0.0.1 - 127.255.255.255
		//255.255.255.255
		//0.0.0.0
		//224.0.0.0 - 255.255.255.255
		//169.254.0.0 - 169.254.255.255
		return ip.IsGlobalUnicast()
	},
}

func NewUID(id string) *UID {
	for i := 0; i < int(UIDTypeMaxLimit); i++ {
		if CheckInput[i](id) {
			return &UID{
				T:  UIDType(i),
				id: id,
			}
		}
	}
	return nil
}

// UID universal ID
type UID struct {
	T  UIDType
	id string
}

func (uid *UID) ToString() string {
	t := strconv.FormatInt(int64(uid.T), 16)
	if len(t) < 2 {
		t = "0" + t
	}
	return t + uid.id
}

func (uid *UID) GetUIDType() string {
	if uid.T < 0 || uid.T >= UIDTypeMaxLimit {
		return "无效的UID类型"
	}
	return IDTypeName[uid.T]
}

func (uid *UID) GetUIDContent() string {
	return uid.id
}

func (uid *UID) Verify() bool {
	if uid.T >= UIDTypeMaxLimit && uid.T < 0 {
		return false
	}
	return verifyUID[uid.T](uid.id)
}

func RetrieveLocalIPV4() ([]string, error) {
	var res []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			res = append(res, ip.String())
		}
	}
	if len(res) == 0 {
		return nil, errors.New("doesn't connect to network")
	}
	return res, nil
}
