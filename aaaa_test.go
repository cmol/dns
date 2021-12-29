package dnsmessage

import (
	"bytes"
	"reflect"
	"testing"

	"inet.af/netaddr"
)

func TestIPv6_Parse(t *testing.T) {
	type args struct {
		buf     []byte
		domains map[int]string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Simple IPv6 address",
			args: args{
				buf: []byte("&\a\xF8\xB0@\v\b\x02\x00\x00\x00\x00\x00\x00 \x11"),
			},
			want: "2607:f8b0:400b:802::2011",
		},
		{
			name: "Not enough data in buffer",
			args: args{
				buf: []byte("&\a\xF8\xB0@\v\b\x02\x00\x00\x00\x00\x00\x00"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := IPv6{}
			if err := ip.Parse(bytes.NewBuffer(tt.args.buf), tt.args.domains); (err != nil) != tt.wantErr {
				t.Errorf("IPv4.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			wantIP, _ := netaddr.ParseIP(tt.want)
			if !reflect.DeepEqual(ip.IP, wantIP) {
				t.Errorf("ParseQuestion() = %v, want %v", ip.IP, wantIP)
			}
		})
	}
}

func TestIPv6_Build(t *testing.T) {
	type args struct {
		domains map[string]int
	}
	tests := []struct {
		name    string
		args    args
		address string
		want    []byte
		wantErr bool
	}{
		{
			name:    "Build simple address",
			address: "2607:f8b0:400b:802::2011",
			want:    []byte("&\a\xF8\xB0@\v\b\x02\x00\x00\x00\x00\x00\x00 \x11"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := &IPv6{}
			ip.IP, _ = netaddr.ParseIP(tt.address)
			buf := new(bytes.Buffer)
			if err := ip.Build(buf, tt.args.domains); (err != nil) != tt.wantErr {
				t.Errorf("IPv4.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.want) {
				t.Errorf("ParseQuestion() = %v, want %v", buf.Bytes(), tt.want)
			}
		})
	}
}
