package snowflake

import (
	"net"
	"os"
)

func WorkerIDFromLocalIP() (uint32, error) {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = os.Getenv("HOSTNAME")
	}

	var ipv4 net.IP
	addresses, err := net.LookupIP(hostname)
	if err != nil {
		return 0, err
	}

	for _, addr := range addresses {
		if ipv4 = addr.To4(); ipv4 != nil {
			break
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
