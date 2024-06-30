package dns

import (
	"bytes"
	"encoding/binary"
)

// Txt implements interface RData
type Txt struct {
	Length uint8
	Data   string
}

// Parse implements TXT parsing for interface RData
func (t *Txt) Parse(buf *bytes.Buffer, _ int, _ *Domains) error {
	if err := binary.Read(buf, binary.BigEndian, &t.Length); err != nil {
		return err
	}
	t.Data = string(buf.Next(int(t.Length)))
	return nil
}

// Build implements TXT building for interface RData
func (t *Txt) Build(buf *bytes.Buffer, _ *Domains) error {
	if err := binary.Write(buf, binary.BigEndian, t.Length); err != nil {
		return err
	}
	buf.WriteString(t.Data)
	return nil
}

// PreBuild step, building name and adding full record
func (t *Txt) PreBuild(_ *Record, _ *Domains) (int, error) {
	return len(t.Data) + 1, nil
}

// TransformName adds service/proto/name fields of server record
func (t *Txt) TransformName(name string) string {
	return name
}
