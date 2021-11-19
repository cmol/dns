package dnsmessage

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParseQuestion(t *testing.T) {
	type args struct {
		buf     []byte
		pointer int
		domains map[int]string
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
				domains: map[int]string{},
			},
			want: Question{
				Domain: "domain.test",
				QType:  A,
				QClass: IN,
			},
			wantErr: false,
		},
		{
			name: "Message missing byte",
			args: args{
				buf:     []byte("\x06domain\x04test\x00\x01\x00\x01"),
				pointer: 0,
				domains: map[int]string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseQuestion(bytes.NewBuffer(tt.args.buf), tt.args.pointer,
				tt.args.domains)
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
		Domain string
		QType  Type
		QClass Class
	}
	type args struct {
		domains map[string]int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "Test simple question",
			fields: fields{
				Domain: "domain.test",
				QType:  A,
				QClass: IN,
			},
			args: args{
				domains: map[string]int{},
			},
			wantErr: false,
			wantBuf: []byte("\x06domain\x04test\x00\x00\x01\x00\x01"),
		},
		{
			name: "Question with missing Type",
			fields: fields{
				Domain: "domain.test",
				QType:  0,
				QClass: IN,
			},
			args: args{
				domains: map[string]int{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Domain: tt.fields.Domain,
				QType:  tt.fields.QType,
				QClass: tt.fields.QClass,
			}
			b := new(bytes.Buffer)
			if err := q.Build(b, tt.args.domains); (err != nil) != tt.wantErr {
				t.Errorf("Question.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(b.Bytes(), tt.wantBuf) {
				t.Errorf("ParseQuestion() = %v, want %v", b.Bytes(), tt.wantBuf)
			}
		})
	}
}
