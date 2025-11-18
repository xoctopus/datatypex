package snowflake

import (
	"net"
)

func WorkerIDFromLocalIP() (uint32, error) {
	var ipv4 net.IP

	addresses, _ := net.InterfaceAddrs()
	for _, addr := range addresses {
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ipv4 = ip.IP.To4(); ipv4 != nil {
				break
			}
		}
	}
	return WorkerIDFromIP(ipv4), nil
}

func WorkerIDFromIP(ipv4 net.IP) uint32 {
	if ipv4 == nil {
		return 0
	}
	ip := ipv4.To4()
	return uint32(ip[2])<<8 + uint32(ip[3])
}
