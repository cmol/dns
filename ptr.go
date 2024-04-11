package dns

import (
	"bytes"
	"errors"
)

// Ptr implements interface RData
type Ptr struct {
	Name  string
	bytes string
}

// Parse implements Ptr parsing for interface RData
func (n *Ptr) Parse(buf *bytes.Buffer, ptr int, domains *Domains) error {
	name, err := ParseName(buf, ptr, domains)
	if err != nil {
		return errors.New("unable to parse Ptr: " + err.Error())
	}
	n.Name = name
	return nil
}

// Build implements Ptr building for interface RData
func (n *Ptr) Build(buf *bytes.Buffer, domains *Domains) error {
	domains.SetBuild(buf.Len(), n.Name)
	buf.WriteString(n.bytes)
	return nil
}

// PreBuild implements Ptr pre building for interface RData
func (n *Ptr) PreBuild(_ *Record, domains *Domains) (int, error) {
	n.bytes = BuildName(n.Name, domains)
	return len(n.bytes), nil
}
