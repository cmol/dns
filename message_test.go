package dns

import (
	"bytes"
	"reflect"
	"testing"

	"net/netip"
)

func TestMessage_ParseHeader(t *testing.T) {
	type fields struct {
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
	type args struct {
		buf []byte
	}
	tests := []struct {
		name       string
		wantFields fields
		args       args
		wantErr    bool
	}{
		{
			name: "Test simple request header",
			wantFields: fields{
				id:      0x199f,
				qr:      false,
				opcode:  0,
				aa:      false,
				tc:      false,
				rd:      true,
				ra:      false,
				rcode:   0,
				qdcount: 1,
				ancount: 0,
				nscount: 0,
				arcount: 0,
			},
			args: args{
				buf: []byte("\x19\x9f\x01\x00\x00\x01\x00\x00\x00\x00\x00\x00"), /*\x06\x64\x6f" +
				"\x6d\x61\x69\x6e\x04\x74\x65\x73\x74\x00\x00\x01\x00\x01\x00\x00\x29" +
				"\x02\x00\x00\x00\x00\x00\x00\x00"),*/
			},
			wantErr: false,
		},
		{
			name: "Test simple response header",
			wantFields: fields{
				id:      0x3028,
				qr:      true,
				opcode:  0,
				aa:      false,
				tc:      false,
				rd:      true,
				ra:      true,
				rcode:   0,
				qdcount: 1,
				ancount: 1,
				nscount: 0,
				arcount: 1,
			},
			args: args{
				buf: []byte("\x30\x28\x81\x80\x00\x01\x00\x01\x00\x00\x00\x01\x07\x63\x6f\x6e" +
					"\x74\x69\x6c\x65\x08\x73\x65\x72\x76\x69\x63\x65\x73\x07\x6d\x6f" +
					"\x7a\x69\x6c\x6c\x61\x03\x63\x6f\x6d\x00\x00\x01\x00\x01\xc0\x0c" +
					"\x00\x01\x00\x01\x00\x00\x00\xb2\x00\x04\x22\x75\xed\xef\x00\x00" +
					"\x29\x10\x00\x00\x00\x00\x00\x00\x00"),
			},
			wantErr: false,
		},
		{
			name:       "Missing fields",
			wantFields: fields{},
			args: args{
				buf: []byte("\x30\x28"),
			},
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantMsg := Message{
				ID:        tt.wantFields.id,
				QR:        tt.wantFields.qr,
				OPCode:    tt.wantFields.opcode,
				AA:        tt.wantFields.aa,
				TC:        tt.wantFields.tc,
				RD:        tt.wantFields.rd,
				RA:        tt.wantFields.ra,
				RCode:     tt.wantFields.rcode,
				qdcount:   tt.wantFields.qdcount,
				ancount:   tt.wantFields.ancount,
				nscount:   tt.wantFields.nscount,
				arcount:   tt.wantFields.arcount,
				Questions: tt.wantFields.questions,
			}
			m := Message{}
			if err := m.ParseHeader(bytes.NewBuffer(tt.args.buf)); (err != nil) != tt.wantErr {
				t.Errorf("Message.ParseHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(m, wantMsg) && !tt.wantErr {
				t.Errorf("ParseHeader() = %v, want %v", m, wantMsg)
			}
		})
	}
}

func TestMessage_BuildHeader(t *testing.T) {
	type fields struct {
		id          uint16
		qr          bool
		opcode      uint8
		aa          bool
		tc          bool
		rd          bool
		ra          bool
		rcode       uint8
		qdcount     uint16
		ancount     uint16
		nscount     uint16
		arcount     uint16
		questions   []Question
		answers     []Record
		nameservers []Record
		additional  []Record
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		want    []byte
	}{
		{
			name: "Test simple request header",
			fields: fields{
				id:      0x199f,
				qr:      false,
				opcode:  0,
				aa:      false,
				tc:      false,
				rd:      true,
				ra:      false,
				rcode:   0,
				qdcount: 1,
				ancount: 0,
				nscount: 0,
				arcount: 0,
			},
			want: []byte("\x19\x9f\x01\x00\x00\x01\x00\x00\x00\x00\x00\x00"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				ID:          tt.fields.id,
				QR:          tt.fields.qr,
				OPCode:      tt.fields.opcode,
				AA:          tt.fields.aa,
				TC:          tt.fields.tc,
				RD:          tt.fields.rd,
				RA:          tt.fields.ra,
				RCode:       tt.fields.rcode,
				qdcount:     tt.fields.qdcount,
				ancount:     tt.fields.ancount,
				nscount:     tt.fields.nscount,
				arcount:     tt.fields.arcount,
				Questions:   tt.fields.questions,
				Answers:     tt.fields.answers,
				Nameservers: tt.fields.nameservers,
				Additional:  tt.fields.additional,
			}
			buf := new(bytes.Buffer)
			if err := m.BuildHeader(buf); (err != nil) != tt.wantErr {
				t.Errorf("Message.BuildHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.want) {
				t.Errorf("Message.BuildHeader() got = %v, want %v", buf.Bytes(), tt.want)
			}
		})
	}
}

func TestParseMessage(t *testing.T) {
	tests := []struct {
		name    string
		buf     []byte
		want    *Message
		wantErr bool
	}{
		{
			name: "Simple response message",
			buf:  []byte("\x00\x1d\x81\x80\x00\x01\x00\x01\x00\x00\x00\x00\x06\x67\x6f\x6c\x61\x6e\x67\x03\x6f\x72\x67\x00\x00\x1c\x00\x01\xc0\x0c\x00\x1c\x00\x01\x00\x00\x01\x2c\x00\x10\x26\x07\xf8\xb0\x40\x0b\x08\x02\x00\x00\x00\x00\x00\x00\x20\x11"),
			want: &Message{
				ID:      0x001d,
				QR:      true,
				RD:      true,
				RA:      true,
				qdcount: 1,
				ancount: 1,
				Questions: []Question{{
					Domain: "golang.org",
					Type:   AAAA,
					Class:  1,
				}},
				Answers: []Record{{
					TTL:    300,
					Class:  1,
					Length: 16,
					Type:   AAAA,
					Name:   "golang.org",
					Data:   &IPv6{Addr: netip.MustParseAddr("2607:f8b0:400b:802::2011")},
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMessage(bytes.NewBuffer(tt.buf))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_Build(t *testing.T) {
	type fields struct {
		id          uint16
		qr          bool
		opcode      uint8
		aa          bool
		tc          bool
		rd          bool
		ra          bool
		rcode       uint8
		qdcount     uint16
		ancount     uint16
		nscount     uint16
		arcount     uint16
		questions   []Question
		answers     []Record
		nameservers []Record
		additional  []Record
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Simple response message",
			want: []byte("\x00\x1d\x81\x80\x00\x01\x00\x01\x00\x00\x00\x00\x06\x67\x6f\x6c\x61\x6e\x67\x03\x6f\x72\x67\x00\x00\x1c\x00\x01\xc0\x0c\x00\x1c\x00\x01\x00\x00\x01\x2c\x00\x10\x26\x07\xf8\xb0\x40\x0b\x08\x02\x00\x00\x00\x00\x00\x00\x20\x11"),
			fields: fields{
				id:      0x001d,
				qr:      true,
				rd:      true,
				ra:      true,
				qdcount: 1,
				ancount: 1,
				questions: []Question{{
					Domain: "golang.org",
					Type:   AAAA,
					Class:  1,
				}},
				answers: []Record{{
					TTL:    300,
					Class:  1,
					Length: 16,
					Type:   AAAA,
					Name:   "golang.org",
					Data:   &IPv6{Addr: netip.MustParseAddr("2607:f8b0:400b:802::2011")},
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				ID:          tt.fields.id,
				QR:          tt.fields.qr,
				OPCode:      tt.fields.opcode,
				AA:          tt.fields.aa,
				TC:          tt.fields.tc,
				RD:          tt.fields.rd,
				RA:          tt.fields.ra,
				RCode:       tt.fields.rcode,
				qdcount:     tt.fields.qdcount,
				ancount:     tt.fields.ancount,
				nscount:     tt.fields.nscount,
				arcount:     tt.fields.arcount,
				Questions:   tt.fields.questions,
				Answers:     tt.fields.answers,
				Nameservers: tt.fields.nameservers,
				Additional:  tt.fields.additional,
			}
			buf := new(bytes.Buffer)
			if err := m.Build(buf, NewDomains()); (err != nil) != tt.wantErr {
				t.Errorf("Message.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.want) {
				t.Errorf("Message.Build() got = %v, want %v", buf.Bytes(), tt.want)
			}
		})
	}
}

func BenchmarkMessageParsing(b *testing.B) {
	buf := []byte("\x00\x1d\x81\x80\x00\x01\x00\x01\x00\x00\x00\x00\x06\x67\x6f\x6c\x61\x6e\x67\x03\x6f\x72\x67\x00\x00\x1c\x00\x01\xc0\x0c\x00\x1c\x00\x01\x00\x00\x01\x2c\x00\x10\x26\x07\xf8\xb0\x40\x0b\x08\x02\x00\x00\x00\x00\x00\x00\x20\x11")
	for i := 0; i < b.N; i++ {
		ParseMessage(bytes.NewBuffer(buf))
	}
}
