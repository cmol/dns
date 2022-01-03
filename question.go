package dnsmessage

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Question struct {
	Domain string
	QType  Type
	QClass Class
}

func ParseQuestion(buf *bytes.Buffer, pointer int, domains *Pointers) (Question, error) {
	q := Question{}
	dom, err := ParseName(buf, pointer, domains)
	if err != nil {
		return Question{}, err
	}
	q.Domain = dom
	if err := binary.Read(buf, binary.BigEndian, &q.QType); err != nil {
		return Question{}, err
	}
	if err := binary.Read(buf, binary.BigEndian, &q.QClass); err != nil {
		return Question{}, err
	}

	return q, nil
}

func (q *Question) Build(buf *bytes.Buffer, domains *Pointers) error {
	if q.Domain == "" || q.QType == 0 || q.QClass == 0 {
		return errors.New("Domain or query type unset")
	}

	BuildName(buf, q.Domain, domains)
	err := binary.Write(buf, binary.BigEndian, q.QType)
	if err != nil {
		return err
	}
	err = binary.Write(buf, binary.BigEndian, q.QClass)
	if err != nil {
		return err
	}
	return nil
}
