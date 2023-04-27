package util

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"net"
	"sort"
	"strings"
)

// GetUniqueIDFromMac get unique id from mac
func GetUniqueIDFromMac() (unique string, err error) {
	addrs := getMacAddrs()
	sort.Strings(addrs)
	addrInOne := strings.Join(addrs, ":")
	hash := md5.New()
	_, err = hash.Write([]byte(addrInOne))
	if err != nil {
		return "", err
	}
	// fmt.Printf("writed: %d\n", n)
	unique = fmt.Sprintf("%x", hash.Sum(nil))
	return
}

func getMacAddrs() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		log.Printf("fail to get net interfaces: %v", err)
		return macAddrs
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}

		macAddrs = append(macAddrs, macAddr)
	}
	return macAddrs
}

func GetIntranetIp() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("ip:", ipnet.IP.String())
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("error to get intranet ip")
}
