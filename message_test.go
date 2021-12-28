package dnsmessage

import (
	"bytes"
	"reflect"
	"testing"
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
				id:        tt.wantFields.id,
				qr:        tt.wantFields.qr,
				opcode:    tt.wantFields.opcode,
				aa:        tt.wantFields.aa,
				tc:        tt.wantFields.tc,
				rd:        tt.wantFields.rd,
				ra:        tt.wantFields.ra,
				rcode:     tt.wantFields.rcode,
				qdcount:   tt.wantFields.qdcount,
				ancount:   tt.wantFields.ancount,
				nscount:   tt.wantFields.nscount,
				arcount:   tt.wantFields.arcount,
				questions: tt.wantFields.questions,
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
