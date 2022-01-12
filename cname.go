package dnsmessage

import (
	"bytes"
	"errors"
)

type CName struct {
	Name string
	buf  []byte
}

func (n *CName) Parse(buf *bytes.Buffer, ptr int, domains *Domains) error {
	name, err := ParseName(buf, ptr, domains)
	if err != nil {
		return errors.New("unable to parse CNAME: " + err.Error())
	}
	n.Name = name
	return nil
}

func (n *CName) Build(buf *bytes.Buffer) error {
	buf.Write(n.buf)
	return nil
}
func (n *CName) PreBuild(domains *Domains) (int, error) {
	buf := new(bytes.Buffer)
	BuildName(buf, n.Name, domains)
	n.buf = buf.Bytes()
	return len(n.buf), nil
}
