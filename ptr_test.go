package dns

import (
	"bytes"
	"reflect"
	"testing"
)

func TestPtr_Parse(t *testing.T) {
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
			n := &Ptr{}
			if err := n.Parse(bytes.NewBuffer(tt.args.buf), tt.args.ptr, tt.args.domains); (err != nil) != tt.wantErr {
				t.Errorf("Ptr.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(n.Name, tt.want) {
				t.Errorf("Ptr.Parse() = %v, want %v", n.Name, tt.want)
			}
		})
	}
}

func TestPtr_Build(t *testing.T) {
	tests := []struct {
		name       string
		recordName string
		want       []byte
		wantLength int
		wantErr    bool
	}{
		{
			name:       "Simple domain",
			recordName: "golang.com",
			want:       []byte("\x06golang\x03com\x00"),
			wantLength: 12,
		},
		{
			name:       "Name with multibyte(3) chars",
			recordName: "ㄶ.golang.com",
			want:       []byte("\x03\xE3\x84\xB6\x06golang\x03com\x00"),
			wantLength: 16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Ptr{
				Name: tt.recordName,
			}
			buf := new(bytes.Buffer)
			if length, err := n.PreBuild(&Record{}, NewDomains()); err != nil || length != tt.wantLength {
				if err != nil {
					t.Errorf("Ptr.PreBuild() error = %v", err)
				} else {
					t.Errorf("Ptr.PreBuild() = %v, want %v", length, tt.wantLength)
				}
				return
			}
			if err := n.Build(buf, NewDomains()); (err != nil) != tt.wantErr {
				t.Errorf("Ptr.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.want) {
				t.Errorf("Ptr.Build() = %v, want %v", buf.Bytes(), tt.want)
			}
		})
	}
}
