package proactor

import (
	"fmt"
	"net"
	"syscall"
)

type Acceptor interface {
	net.Listener
}

type acceptor struct {
	family   int
	type_    int
	protocol int
	fd       syscall.Handle
}

func Listen(domain int, typ int, proto int, addr string, port uint16) (Acceptor, error) {
	fd, err := syscall.Socket(domain, typ, proto)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			syscall.Closesocket(fd)
		}
	}()

	object := &acceptor{
		family:   domain,
		type_:    typ,
		protocol: proto,
		fd:       fd,
	}

	if err = object.bind(addr, port); err != nil {
		return nil, err
	}

	return object, nil
}

func (a *acceptor) Accept() (net.Conn, error) {
	nfd, sa, err := syscall.Accept(a.fd)
	if err != nil {
		return nil, err
	}

	link := &linker{fd: nfd, localAddr: sa}
	return link, nil
}

func (a *acceptor) Close() error {
	return syscall.Closesocket(a.fd)
}

// Addr returns the listener's network address.
func (a *acceptor) Addr() net.Addr {
	return nil
}

func (a *acceptor) bind(addr string, port uint16) error {
	ip := net.ParseIP(addr)
	if ip == nil {
		return fmt.Errorf("invalid IP address %q", addr)
	}

	if ip.To4() == nil {
		ip = net.IPv4zero
	}

	sockAddr := &syscall.SockaddrInet4{Port: int(port)}
	copy(sockAddr.Addr[:], ip.To4())

	if err := syscall.Bind(a.fd, sockAddr); err != nil {
		return err
	}

	return syscall.Listen(a.fd, syscall.SOMAXCONN)
}
