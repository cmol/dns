package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Txt implements interface RData
type Txt struct {
  Length uint16
	Data   []string
}

// Parse implements TXT parsing for interface RData
func (t *Txt) Parse(buf *bytes.Buffer, _ int, _ *Domains) error {
  var read uint = 0
	for {
		len, err := buf.ReadByte()
		if err != nil {
			return err
		}

		if uint(len)+read > uint(t.Length) {
      return fmt.Errorf(
        "txt record part length too long: %d > %d\nstrings read: %+v",
        uint(len)+read, t.Length, t.Data,
      )
		}

		t.Data = append(t.Data, string(buf.Next(int(len))))
		read = read + uint(len) + 1
		if read == uint(t.Length) {
			break
		}
	}
	return nil
}

// Build implements TXT building for interface RData
func (t *Txt) Build(buf *bytes.Buffer, _ *Domains) error {
	for _, part := range t.Data {
    partLen := uint8(len(part))
		if err := binary.Write(buf, binary.BigEndian, partLen); err != nil {
			return err
		}
		buf.WriteString(part)
	}
	return nil
}

// PreBuild step, building name and adding full record
func (t *Txt) PreBuild(_ *Record, _ *Domains) (int, error) {
	writeLength := 0
	for _, s := range t.Data {
		writeLength = writeLength + len(s) + 1 // Add 1 for the length indicator
	}

	return writeLength, nil
}

// TransformName adds service/proto/name fields of server record
func (t *Txt) TransformName(name string) string {
	return name
}
