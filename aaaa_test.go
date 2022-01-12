package dnsmessage

import (
	"bytes"
	"reflect"
	"testing"

	"inet.af/netaddr"
)

func TestIPv6_Parse(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		args    args
		want    IPv6
		wantErr bool
	}{
		{
			name: "Simple IPv6 address",
			args: args{buf: []byte("&\a\xF8\xB0@\v\b\x02\x00\x00\x00\x00\x00\x00 \x11")},
			want: IPv6{netaddr.MustParseIP("2607:f8b0:400b:802::2011")},
		},
		{
			name:    "Not enough data in buffer",
			args:    args{buf: []byte("&\a\xF8\xB0@\v\b\x02\x00\x00\x00\x00\x00\x00")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := IPv6{}
			if err := ip.Parse(bytes.NewBuffer(tt.args.buf), 0, &Domains{}); (err != nil) != tt.wantErr {
				t.Errorf("IPv6.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(ip.IP, tt.want.IP) {
				t.Errorf("IPv6.Parse() = %v, want %v", ip.IP, tt.want.IP)
			}
		})
	}
}

func TestIPv6_Build(t *testing.T) {
	tests := []struct {
		name    string
		address netaddr.IP
		want    []byte
		wantErr bool
	}{
		{
			name:    "Build simple address",
			address: netaddr.MustParseIP("2607:f8b0:400b:802::2011"),
			want:    []byte("&\a\xF8\xB0@\v\b\x02\x00\x00\x00\x00\x00\x00 \x11"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := &IPv6{IP: tt.address}
			buf := new(bytes.Buffer)
			if err := ip.Build(buf); (err != nil) != tt.wantErr {
				t.Errorf("IPv4.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.want) {
				t.Errorf("IPv6.Build() = %v, want %v", buf.Bytes(), tt.want)
			}
		})
	}
}
