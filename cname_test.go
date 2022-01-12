package dnsmessage

import (
	"bytes"
	"reflect"
	"testing"
)

func TestCName_Parse(t *testing.T) {
	type args struct {
		buf     []byte
		ptr     int
		domains *Domains
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Simple domain name",
			args: args{
				buf:     []byte("\x06golang\x03com\x00"),
				domains: NewDomains(),
			},
			want: "golang.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &CName{}
			if err := n.Parse(bytes.NewBuffer(tt.args.buf), tt.args.ptr, tt.args.domains); (err != nil) != tt.wantErr {
				t.Errorf("CName.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(n.Name, tt.want) {
				t.Errorf("CName.Parse() = %v, want %v", n.Name, tt.want)
			}
		})
	}
}

func TestCName_Build(t *testing.T) {
	tests := []struct {
		name       string
		cname      string
		want       []byte
		wantLength int
		wantErr    bool
	}{
		{
			name:       "Simple domain",
			cname:      "golang.com",
			want:       []byte("\x06golang\x03com\x00"),
			wantLength: 12,
		},
		{
			name:       "Name with multibyte(3) chars",
			cname:      "ㄶ.golang.com",
			want:       []byte("\x03\xE3\x84\xB6\x06golang\x03com\x00"),
			wantLength: 16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &CName{
				Name: tt.cname,
			}
			buf := new(bytes.Buffer)
			if length, err := n.PreBuild(NewDomains()); err != nil || length != tt.wantLength {
				if err != nil {
					t.Errorf("CName.PreBuild() error = %v", err)
				} else {
					t.Errorf("CName.PreBuild() = %v, want %v", length, tt.wantLength)
				}
				return
			}
			if err := n.Build(buf); (err != nil) != tt.wantErr {
				t.Errorf("CName.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.want) {
				t.Errorf("CName.Build() = %v, want %v", buf.Bytes(), tt.want)
			}
		})
	}
}
