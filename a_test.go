package dnsmessage

import (
	"bytes"
	"reflect"
	"testing"

	"inet.af/netaddr"
)

func TestIPv4_Parse(t *testing.T) {
	type args struct {
		buf     []byte
		domains map[int]string
	}
	tests := []struct {
		name    string
		args    args
		want    IPv4
		wantErr bool
	}{
		{
			name: "Simple IPv4 address",
			args: args{
				buf: []byte("\xC0\xA8\x16\xBC"),
			},
			want: IPv4{
				IP: netaddr.IPv4(192, 168, 22, 188),
			},
		},
		{
			name: "Not enough data in buffer",
			args: args{
				buf: []byte("\xC0\xA8\x16"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := IPv4{}
			if err := ip.Parse(bytes.NewBuffer(tt.args.buf), tt.args.domains); (err != nil) != tt.wantErr {
				t.Errorf("IPv4.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(ip, tt.want) {
				t.Errorf("ParseQuestion() = %v, want %v", ip, tt.want)
			}
		})
	}
}

func TestIPv4_Build(t *testing.T) {
	type fields struct {
		IP netaddr.IP
	}
	type args struct {
		domains map[string]int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Build simple address",
			fields: fields{
				IP: netaddr.IPv4(192, 168, 22, 188),
			},
			want: []byte("\xC0\xA8\x16\xBC"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := &IPv4{
				IP: tt.fields.IP,
			}
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
