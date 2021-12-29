package dnsmessage

import (
	"bytes"
	"reflect"
	"testing"

	"inet.af/netaddr"
)

func TestParseRecord(t *testing.T) {
	type args struct {
		buf     []byte
		pointer int
		domains map[int]string
	}
	tests := []struct {
		name    string
		args    args
		want    Record
		wantErr bool
	}{
		{
			name: "Simple AAAA record",
			args: args{
				buf:     []byte("\x06golang\x03com\x00\x00\x1c\x00\x01\x00\x00\x01\x2c\x00\x10\x26\x07\xf8\xb0\x40\x0b\x08\x02\x00\x00\x00\x00\x00\x00\x20\x11"),
				domains: make(map[int]string),
			},
			want: Record{
				Name:        "golang.com",
				RType:       AAAA,
				Class:       1,
				TTL:         300,
				RDataLength: 16,
				Data: &IPv6{
					IP: netaddr.MustParseIP("2607:f8b0:400b:802::2011"),
				},
			},
		},
		{
			name: "Simple A record",
			args: args{
				buf:     []byte("\x06golang\x03com\x00\x00\x01\x00\x01\x00\x00\x01\x2c\x00\x04\x8e\xfb\x29\x51"),
				domains: make(map[int]string),
			},
			want: Record{
				Name:        "golang.com",
				RType:       A,
				Class:       1,
				TTL:         300,
				RDataLength: 4,
				Data: &IPv4{
					IP: netaddr.MustParseIP("142.251.41.81"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRecord(bytes.NewBuffer(tt.args.buf), tt.args.pointer, tt.args.domains)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}
