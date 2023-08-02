package dns

import (
	"bytes"
	"errors"
)

// CName implements interface RData
type CName struct {
	Name  string
	bytes string
}

// Parse implements CNAME parsing for interface RData
func (n *CName) Parse(buf *bytes.Buffer, ptr int, domains *Domains) error {
	name, err := ParseName(buf, ptr, domains)
	if err != nil {
		return errors.New("unable to parse CNAME: " + err.Error())
	}
	n.Name = name
	return nil
}

// Build implements CNAME building for interface RData
func (n *CName) Build(buf *bytes.Buffer, domains *Domains) error {
	domains.SetBuild(buf.Len(), n.Name)
	buf.WriteString(n.bytes)
	return nil
}

// PreBuild implements CNAME pre building for interface RData
func (n *CName) PreBuild(r *Record, domains *Domains) (int, error) {
	n.bytes = BuildName(n.Name, domains)
	return len(n.bytes), nil
}
