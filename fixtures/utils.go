package fixtures

import (
	"net"
	"os"
)

// Get local IPv4
func GetLocalIPv4() (ipv4 net.IP, err error) {
	host, err := os.Hostname()
	if err != nil {
		return
	}
	addrs, err := net.LookupIP(host)
	if err != nil {
		return
	}
	for _, addr := range addrs {
		if ipv4 = addr.To4(); ipv4 != nil {
			return
		}
	}
	return
}

func GetOpenPort() (p int, err error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}
	defer l.Close()
	p = l.Addr().(*net.TCPAddr).Port
	return
}

func GetOpenPorts(n int) (p []int, err error) {
	for i := 0; i < n; i++ {
		var addr *net.TCPAddr
		var listener net.Listener
		if addr, err = net.ResolveTCPAddr("tcp", "localhost:0"); err != nil {
			return
		}
		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			return
		}
		defer listener.Close()
		p = append(p, listener.Addr().(*net.TCPAddr).Port)
	}
	return
}
