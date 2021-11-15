package dnsmessage

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

const NAME_POINTER = 0xc0
const POINTER_MASK = 0x3fff

func ParseName(bytes []byte, pointer int, domains map[int]string) (string, error) {
	if len(bytes) < pointer {
		return "", errors.New("Unable to parse name, pointer oob")
	}
	length := int(bytes[pointer])

	// Check if name is a pointer to an earlier refereced domain
	if length&NAME_POINTER == NAME_POINTER {
		n, ok := domains[(length<<8|int(bytes[pointer+1]))&POINTER_MASK]
		if !ok {
			return "", errors.New("Name pointer points to nothing")
		}
		return n, nil
	}

	strname, err := parseNameBytes(bytes, pointer, length)
	if err != nil {
		return "", err
	}
	domains[pointer] = strname
	return strname, nil
}

func parseNameBytes(bytes []byte, pointer, length int) (string, error) {
	bytesLength := len(bytes)
	var name strings.Builder
	for length > 0 {
		if bytesLength < pointer+length+1 {
			return "", errors.New("Unable to parse name, subdomain oob")
		}
		name.Write(bytes[pointer+1 : pointer+length+1])
		pointer += length + 1
		length = int(bytes[pointer])
		if length > 0 {
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
