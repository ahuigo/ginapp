package nettool

import (
	"fmt"
	"log"
	"net"
)

func GetLocalIP() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("get local ip:", err.Error())
		return ""
	}

	for _, addr := range addresses {
		var ip net.IP
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.IsLoopback() {
			continue
		}
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		// process IP address
		if ip.To4() != nil {
			return ip.String()
		}
	}
	return ""
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
