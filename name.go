package dnsmessage

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

const NAME_POINTER = 0xc0
const POINTER_MASK = 0x3fff

func ParseName(buf *bytes.Buffer, ptr int, domains map[int]string) (string, error) {
	length, err := buf.ReadByte()
	if err != nil {
		return "", err
	}

	// Check if name is a pointer to an earlier refereced domain
	if length&NAME_POINTER == NAME_POINTER {
		l2, err := buf.ReadByte()
		if err != nil {
			return "", err
		}
		n, ok := domains[(int(length)<<8|int(l2))&POINTER_MASK]
		if !ok {
			return "", errors.New("Name pointer points to nothing")
		}
		return n, nil
	}

	strname, err := parseNameBytes(buf, length)
	if err != nil {
		return "", err
	}
	domains[ptr] = strname
	return strname, nil
}

func parseNameBytes(buf *bytes.Buffer, length uint8) (string, error) {
	var name strings.Builder
	err := *new(error)
	for length > 0 {
		n := buf.Next(int(length))
		if len(n) < int(length) {
			return "", errors.New("Unable to parse name, subdomain oob")
		}
		name.Write(n)
		length, err = buf.ReadByte()
		if err != nil {
			return "", err
		}
		if length != 0 {
			name.WriteString(".")
		}
	}
	return name.String(), nil
}

func BuildName(buf *bytes.Buffer, name string, domains map[string]int) int {
	if n, ok := domains[name]; ok {
		binary.Write(buf, binary.BigEndian, uint16(n|NAME_POINTER<<8))
		return 2
	}

	written := 0
	for _, dom := range strings.Split(name, ".") {
		buf.WriteByte(uint8(len(dom)))
		buf.WriteString(dom)
		written += len(dom) + 1
	}
	buf.WriteByte('\x00')
	return written + 1
}
