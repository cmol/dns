package dnsmessage

import (
	"bytes"
	"encoding/binary"
)

type Question struct {
	Domain string
	QType  Type
	QClass Class
}

func ParseQuestion(buf *bytes.Buffer, pointer int, domains map[int]string) (Question, error) {
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
