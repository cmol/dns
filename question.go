package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Question struct {
	Domain string
	Type   Type
	Class  Class
}

func ParseQuestion(buf *bytes.Buffer, pointer int, domains *Domains) (Question, error) {
	q := Question{}
	dom, err := ParseName(buf, pointer, domains)
	if err != nil {
		return Question{}, err
	}
	q.Domain = dom
	if err := binary.Read(buf, binary.BigEndian, &q.Type); err != nil {
		return Question{}, err
	}
	if err := binary.Read(buf, binary.BigEndian, &q.Class); err != nil {
		return Question{}, err
	}

	return q, nil
}

func (q *Question) Build(buf *bytes.Buffer, domains *Domains) error {
	if q.Domain == "" || q.Type == 0 || q.Class == 0 {
		return errors.New("Domain or query type unset")
	}

	name := BuildName(q.Domain, domains)
	domains.SetBuild(buf.Len(), q.Domain)
	buf.WriteString(name)
	err := binary.Write(buf, binary.BigEndian, q.Type)
	if err != nil {
		return err
	}
	err = binary.Write(buf, binary.BigEndian, q.Class)
	if err != nil {
		return err
	}
	return nil
}
