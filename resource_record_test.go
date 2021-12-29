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
				Data:        &IPv6{netaddr.MustParseIP("2607:f8b0:400b:802::2011")},
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
				Data:        &IPv4{netaddr.MustParseIP("142.251.41.81")},
			},
		},
		{
			name: "Simple A record",
			args: args{
				buf:     []byte("\x06golang\x03com\x00\x00\x77\x00\x01\x00\x00\x01\x2c\x00\x04\x8e\xfb\x29\x51"),
				domains: make(map[int]string),
			},
			wantErr: true,
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

func TestRecord_Build(t *testing.T) {
	type fields struct {
		TTL         uint32
		Class       uint16
		RDataLength uint16
		RType       Type
		Name        string
		Data        RData
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
			name: "Build Simple AAAA record",
			fields: fields{
				TTL:         300,
				Class:       1,
				RDataLength: 16,
				RType:       AAAA,
				Name:        "golang.com",
				Data:        &IPv6{netaddr.MustParseIP("2607:f8b0:400b:802::2011")},
			},
			want: []byte("\x06golang\x03com\x00\x00\x1c\x00\x01\x00\x00\x01\x2c\x00\x10\x26\x07\xf8\xb0\x40\x0b\x08\x02\x00\x00\x00\x00\x00\x00\x20\x11"),
		},
		{
			name: "Build Simple A record",
			fields: fields{
				TTL:         300,
				Class:       1,
				RDataLength: 4,
				RType:       A,
				Name:        "golang.com",
				Data:        &IPv4{netaddr.MustParseIP("142.251.41.81")},
			},
			want: []byte("\x06golang\x03com\x00\x00\x01\x00\x01\x00\x00\x01\x2c\x00\x04\x8e\xfb\x29\x51"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Record{
				TTL:         tt.fields.TTL,
				Class:       tt.fields.Class,
				RDataLength: tt.fields.RDataLength,
				RType:       tt.fields.RType,
				Name:        tt.fields.Name,
				Data:        tt.fields.Data,
			}
			buf := new(bytes.Buffer)
			if err := r.Build(buf, tt.args.domains); (err != nil) != tt.wantErr {
				t.Errorf("Record.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.want) {
				t.Errorf("Record.Build() = %v, want %v", buf.Bytes(), tt.want)
			}
		})
	}
}
