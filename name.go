package dnsmessage

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const NamePointer = 0xc0
const PointerMask = 0x3fff

func ParseName(buf *bytes.Buffer, ptr int, domains *Domains) (string, error) {
	var name bytes.Buffer
	length, err := buf.ReadByte()
	if err != nil {
		return "", err
	}

	// root domain
	if length == 0 {
		return "", nil
	}

	for {
		// Check if name is a pointer to an earlier refereced domain
		if length&NamePointer == NamePointer {
			l2, err := buf.ReadByte()
			if err != nil {
				return "", err
			}
			n, ok := domains.GetParse((int(length)<<8 | int(l2)) & PointerMask)
			if !ok {
				return "", errors.New("Name pointer points to nothing")
			}
			name.WriteString(n)
			return name.String(), nil
		}
		n := buf.Next(int(length))
		name.Write(n)
		length, err = buf.ReadByte()
		if err != nil {
			return "", err
		}
		if length != 0 {
			name.WriteString(".")
		} else {
			break
		}
	}
	domains.SetParse(ptr, name.String())
	return name.String(), nil
}

func BuildName(name string, domains *Domains) string {
	// root domain
	if len(name) == 0 {
		return "\x00"
	}

	var buf bytes.Buffer
	written := 0
	name = name + "."
	for i, c := range name {
		if c == '.' {
			if n, ok := domains.GetBuild(name[written : len(name)-1]); ok {
				binary.Write(&buf, binary.BigEndian, uint16(n|NamePointer<<8))
				return buf.String()
			}
			buf.WriteByte(uint8(i - written))
			buf.WriteString(name[written:i])
			written = written + (i - written) + 1
		}
	}
	buf.WriteByte('\x00')
	return buf.String()
}
