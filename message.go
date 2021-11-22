package dnsmessage

import (
	"bytes"
	"encoding/binary"
)

const HDR_LENGTH = 12

type Message struct {
	id        uint16
	qr        bool
	opcode    uint8
	aa        bool
	tc        bool
	rd        bool
	ra        bool
	rcode     uint8
	qdcount   uint16
	ancount   uint16
	nscount   uint16
	arcount   uint16
	questions []Question
}

func (m *Message) ParseHeader(buf *bytes.Buffer) error {
	var opts uint16
	fields := []interface{}{m.id, opts, m.qdcount, m.ancount, m.nscount, m.arcount}
	for _, f := range fields {
		err := binary.Read(buf, binary.BigEndian, &f)
		if err != nil {
			return err
		}
	}
	m.parseOpts(opts)
	return nil
}

func (m *Message) parseOpts(opts uint16) {
	m.rcode = uint8(opts & 0xf)
	m.opcode = uint8((opts >> 11) & 0xf)
	m.ra = opts&OPT_RA == OPT_RA
	m.rd = opts&OPT_RD == OPT_RD
	m.tc = opts&OPT_TC == OPT_TC
	m.aa = opts&OPT_AA == OPT_AA
	m.qr = opts&OPT_QR == OPT_QR
}

func ParseMessage(buf *bytes.Buffer) (Message, error) {
	m := new(Message)
	var err error
	err = m.ParseHeader(buf)
	if err != nil {
		return Message{}, err
	}
	ptr := HDR_LENGTH
	bufLen := len(buf)
	domains := map[int]string{}
	for i := 0; i < m.qdcount; i++ {
		q, err := ParseQuestion(buf, ptr, domains)
		if err != nil {
			return Message{}, err
		}
		m.questions = append(m.questions, q)
		ptr = ptr + bufLen - len(buf)
		bufLen = len(buf)
	}
}
