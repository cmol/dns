package dnsmessage

import (
	"errors"
	"strings"
)

const NAME_POINTER = 0xc0
const POINTER_MASK = 0x3fff

func (m *Message) ParseName(bytes []byte, pointer int) (string, error) {
	var name strings.Builder
	initialPointer := pointer
	bytesLength := len(bytes)

	if bytesLength < pointer {
		return "", errors.New("Unable to parse name, pointer oob")
	}
	length := int(bytes[pointer])

	// check pointer
	if length&NAME_POINTER == NAME_POINTER {
		if n, ok := m.pointers[(length<<8|int(bytes[pointer+1]))&POINTER_MASK]; ok {
			return n, nil
		} else {
			return "", errors.New("Name pointer points to nothing")
		}
	}

	for length > 0 {
		if bytesLength < pointer+length+1 {
			return "", errors.New("Unable to parse name, subdomain oob")
		}
		name.Write(bytes[pointer+1 : pointer+length+1])

		pointer += length + 1
		if bytesLength < pointer {
			return "", errors.New("Unable to parse name, length oob")
		}
		length = int(bytes[pointer])
		if length > 0 {
			name.WriteString(".")
		}
	}
	strname := name.String()
	m.pointers[initialPointer] = strname
	return strname, nil
}
