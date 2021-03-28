package proactor

import (
	"fmt"
	"net"
	"syscall"
	"time"
)

type linker struct {
	fd        syscall.Handle
	localAddr syscall.Sockaddr
}

func Dial(domain int, typ int, proto int, addr string, port uint16) (net.Conn, error) {
	fd, err := syscall.Socket(domain, typ, proto)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			syscall.Closesocket(fd)
		}
	}()

	ip := net.ParseIP(addr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address %q", addr)
	}

	if ip.To4() == nil {
		ip = net.IPv4zero
	}

	sockAddr := &syscall.SockaddrInet4{Port: int(port)}
	copy(sockAddr.Addr[:], ip.To4())

	if err = syscall.Connect(fd, sockAddr); err != nil {
		return nil, err
	}

	return NewLink(fd, sockAddr), err
}

func NewLink(fd syscall.Handle, addr syscall.Sockaddr) net.Conn {
	return &linker{fd: fd, localAddr: addr}
}

func (c *linker) Read(b []byte) (n int, err error) {
	return 0, err
}
func (c *linker) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (c *linker) Close() error {
	return nil
}

func (c *linker) LocalAddr() net.Addr {
	return nil
}

func (c *linker) RemoteAddr() net.Addr {
	return nil
}

func (c *linker) SetDeadline(t time.Time) error {
	return nil
}

func (c *linker) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *linker) SetWriteDeadline(t time.Time) error {
	return nil
}
