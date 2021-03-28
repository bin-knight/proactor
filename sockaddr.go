package proactor

import (
	"fmt"
	"net"
	"syscall"
)

func AddrToIp4(addr string, port uint16) (*syscall.SockaddrInet4, error) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address %q", addr)
	}

	if ip.To4() == nil {
		ip = net.IPv4zero
	}

	sa := &syscall.SockaddrInet4{Port: int(port)}
	copy(sa.Addr[:], ip.To4())
	return sa, nil
}
