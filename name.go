package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	// NamePointer indication byte
	NamePointer = 0xc0

	// PointerMask is the reverse of the indication byte for 16 bits
	PointerMask = 0x3fff
)

// ParseName returns a name given a pointer
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

	newDomain := false

	for {
		// Check if name is a pointer to an earlier refereced domain
		if length&NamePointer == NamePointer {
			l2, err := buf.ReadByte()
			if err != nil {
				return "", err
			}
			getPointer := (int(length)<<8 | int(l2)) & PointerMask
			n, ok := domains.GetParse(getPointer)
			if !ok {
				return "", fmt.Errorf("name pointer %d points to nothing, full map \n%+v",
					getPointer, domains.parsePtr)
			}
			name.WriteString(n)
			if newDomain {
				domains.SetParse(ptr, name.String())
			}
			return name.String(), nil
		}
		newDomain = true
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

// BuildName returns a dns encoded name with pointers if possible
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
				binary.Write(&buf, binary.BigEndian, uint16(NamePointer<<8|n&0xff))
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
