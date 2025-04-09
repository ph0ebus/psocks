package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	_ "log/slog"
	"net"
	"strconv"
)

// METHOD
const (
	NOAUTH byte = iota
	GSSAPI
	AUTH
)

// AYTP
const (
	IPV4 byte = iota + 1
	_
	HOSTNAME
	IPV6
)

// REP
const (
	SUCCESS byte = iota
	FAIL
)

func Socks5Auth(client net.Conn) error {
	buf := make([]byte, 256)
	length, err := io.ReadFull(client, buf[:2])
	if length != 2 {
		return errors.New("failed to read header: " + err.Error())
	}

	// fmt.Printf("%x\n", buf[:2])

	ver, nMethods := buf[0], buf[1]
	if ver != 5 {
		return fmt.Errorf("nnsupported version: %d", ver)
	}

	length, err = io.ReadFull(client, buf[:nMethods])
	if length != int(nMethods) {
		return errors.New("failed to read methods: " + err.Error())
	}
	// fmt.Printf("%x\n", buf[:nMethods])

	_, err = client.Write([]byte{0x05, NOAUTH})
	if err != nil {
		return errors.New("failed to write data: " + err.Error())
	}

	return nil
}

func Socks5Connect(client net.Conn) (net.Conn, error) {
	buf := make([]byte, 256)
	length, err := io.ReadFull(client, buf[:4])
	if length != 4 {
		return nil, errors.New("failed to read header: " + err.Error())
	}

	ver, cmd, _, atyp := buf[0], buf[1], buf[2], buf[3]
	if ver != 5 {
		return nil, fmt.Errorf("unsupported version: %d", ver)
	}

	if cmd != 1 {
		return nil, fmt.Errorf("unsupported cmd: %d", cmd)
	}

	addr := ""

	switch atyp {
	case IPV4:
		length, err = io.ReadFull(client, buf[:4])
		if length != 4 {
			return nil, errors.New("invalid IPV4: " + err.Error())
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])

	case HOSTNAME:
		length, err = io.ReadFull(client, buf[:1])
		if length != 1 {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addrLen := int(buf[0])

		length, err = io.ReadFull(client, buf[:addrLen])
		if length != addrLen {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addr = string(buf[:addrLen])

	default:
		return nil, fmt.Errorf("unsupported atyp: %d", atyp)
	}

	length, err = io.ReadFull(client, buf[:2])
	if length != 2 {
		return nil, errors.New("failed to read port: " + err.Error())
	}
	port := binary.BigEndian.Uint16(buf[:2])

	destAddrPort := net.JoinHostPort(addr, strconv.Itoa(int(port)))
	dest, err := net.Dial("tcp", destAddrPort)
	if err != nil {
		return nil, errors.New("dial dst: " + err.Error())
	}

	_, err = client.Write([]byte{0x05, SUCCESS, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		dest.Close()
		return nil, errors.New("failed to write rsp: " + err.Error())
	}

	return dest, nil

}

func Socks5Forward(client net.Conn, target net.Conn) {

	forward := func(src, dest net.Conn) {
		defer src.Close()
		defer dest.Close()
		io.Copy(src, dest)
	}
	go forward(client, target)
	go forward(target, client)
}
