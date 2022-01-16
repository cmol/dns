package dnsmessage

import (
	"bytes"
	"errors"
)

type CName struct {
	Name  string
	bytes string
}

func (n *CName) Parse(buf *bytes.Buffer, ptr int, domains *Domains) error {
	name, err := ParseName(buf, ptr, domains)
	if err != nil {
		return errors.New("unable to parse CNAME: " + err.Error())
	}
	n.Name = name
	return nil
}

func (n *CName) Build(buf *bytes.Buffer, domains *Domains) error {
	domains.SetBuild(buf.Len(), n.Name)
	buf.WriteString(n.bytes)
	return nil
}

func (n *CName) PreBuild(domains *Domains) (int, error) {
	n.bytes = BuildName(n.Name, domains)
	return len(n.bytes), nil
}
