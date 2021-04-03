package proactor

import (
	"net"
	"syscall"
)

type Addr interface {
	net.Addr
	Sockaddr() (syscall.Sockaddr, error)
}

type TCPAddr net.TCPAddr
type UDPAddr net.UDPAddr

func (taddr *TCPAddr) Sockaddr() (syscall.Sockaddr, error) {
	return ipToSockaddr(addrFamily((*net.TCPAddr)(taddr)), taddr.IP, taddr.Port, taddr.Zone)
}

func (uaddr *UDPAddr) Sockaddr() (syscall.Sockaddr, error) {
	return ipToSockaddr(addrFamily((*net.UDPAddr)(uaddr)), uaddr.IP, uaddr.Port, uaddr.Zone)
}

func isIPv4(addr net.Addr) bool {
	switch addr := addr.(type) {
	case *net.TCPAddr:
		return addr.IP.To4() != nil
	case *net.UDPAddr:
		return addr.IP.To4() != nil
	case *net.IPAddr:
		return addr.IP.To4() != nil
	}
	return false
}

func addrFamily(addr net.Addr) (family int) {
	if isIPv4(addr) {
		return syscall.AF_INET
	}

	return syscall.AF_INET6
}

func ipToSockaddr(family int, ip net.IP, port int, zone string) (syscall.Sockaddr, error) {
	switch family {
	case syscall.AF_INET:
		if len(ip) == 0 {
			ip = net.IPv4zero
		}
		ip4 := ip.To4()
		if ip4 == nil {
			return nil, &net.AddrError{Err: "non-IPv4 address", Addr: ip.String()}
		}
		sa := &syscall.SockaddrInet4{Port: port}
		copy(sa.Addr[:], ip4)
		return sa, nil
	case syscall.AF_INET6:
		if len(ip) == 0 || ip.Equal(net.IPv4zero) {
			ip = net.IPv6zero
		}
		ip6 := ip.To16()
		if ip6 == nil {
			return nil, &net.AddrError{Err: "non-IPv6 address", Addr: ip.String()}
		}
		itf, err := net.InterfaceByName(zone)
		if err != nil {
			return nil, err
		}

		sa := &syscall.SockaddrInet6{Port: port, ZoneId: uint32(itf.Index)}
		copy(sa.Addr[:], ip6)
		return sa, nil
	}

	return nil, &net.AddrError{Err: "invalid address family", Addr: ip.String()}
}
