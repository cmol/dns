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
