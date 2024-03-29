package dns

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParseQuestion(t *testing.T) {
	type args struct {
		buf     []byte
		pointer int
	}
	tests := []struct {
		name    string
		args    args
		want    Question
		wantErr bool
	}{
		{
			name: "Simple single domain.test IN A question",
			args: args{
				buf:     []byte("\x06domain\x04test\x00\x00\x01\x00\x01"),
				pointer: 0,
			},
			want: Question{
				Domain:          "domain.test",
				Type:            A,
				Class:           IN,
				UnicastResponse: false,
			},
			wantErr: false,
		},
		{
			name: "Simple mDNS single domain.test IN A question",
			args: args{
				buf:     []byte("\x06domain\x04test\x00\x00\x01\x80\x01"),
				pointer: 0,
			},
			want: Question{
				Domain:          "domain.test",
				Type:            A,
				Class:           IN,
				UnicastResponse: true,
			},
			wantErr: false,
		},
		{
			name: "Message missing byte",
			args: args{
				buf:     []byte("\x06domain\x04test\x00\x01\x00\x01"),
				pointer: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseQuestion(bytes.NewBuffer(tt.args.buf), tt.args.pointer,
				NewDomains())
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseQuestion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseQuestion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestion_Build(t *testing.T) {
	type fields struct {
		Domain          string
		Type            Type
		Class           Class
		UnicastResponse bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "Test simple question",
			fields: fields{
				Domain:          "domain.test",
				Type:            A,
				Class:           IN,
				UnicastResponse: false,
			},
			wantErr: false,
			wantBuf: []byte("\x06domain\x04test\x00\x00\x01\x00\x01"),
		},
		{
			name: "Test simple mDNS with UnicastResponse question",
			fields: fields{
				Domain:          "domain.test",
				Type:            A,
				Class:           IN,
				UnicastResponse: true,
			},
			wantErr: false,
			wantBuf: []byte("\x06domain\x04test\x00\x00\x01\x80\x01"),
		},
		{
			name: "Question with missing Type",
			fields: fields{
				Domain: "domain.test",
				Type:   0,
				Class:  IN,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Domain:          tt.fields.Domain,
				Type:            tt.fields.Type,
				Class:           tt.fields.Class,
				UnicastResponse: tt.fields.UnicastResponse,
			}
			b := new(bytes.Buffer)
			if err := q.Build(b, NewDomains()); (err != nil) != tt.wantErr {
				t.Errorf("Question.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(b.Bytes(), tt.wantBuf) {
				t.Errorf("ParseQuestion() = %v, want %v", b.Bytes(), tt.wantBuf)
			}
		})
	}
}

func BenchmarkQuestionParsing(b *testing.B) {
	buf := []byte("\x06domain\x04test\x00\x01\x00\x01")
	for i := 0; i < b.N; i++ {
		ParseQuestion(bytes.NewBuffer(buf), 0, NewDomains())
	}
}
