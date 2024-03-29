package dns

import (
	"bytes"
	"encoding/hex"
	"net/netip"
	"reflect"
	"testing"
)

func TestParseRecord(t *testing.T) {
	type args struct {
		buf     []byte
		pointer int
		domains *Domains
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
				domains: NewDomains(),
			},
			want: Record{
				Name:   "golang.com",
				Type:   AAAA,
				Class:  1,
				TTL:    300,
				Length: 16,
				Data:   &IPv6{netip.MustParseAddr("2607:f8b0:400b:802::2011")},
			},
		},
		{
			name: "Simple AAAA record flush mDNS cache",
			args: args{
				buf:     []byte("\x06golang\x03com\x00\x00\x1c\x80\x01\x00\x00\x01\x2c\x00\x10\x26\x07\xf8\xb0\x40\x0b\x08\x02\x00\x00\x00\x00\x00\x00\x20\x11"),
				domains: NewDomains(),
			},
			want: Record{
				Name:       "golang.com",
				Type:       AAAA,
				CacheFlush: true,
				Class:      1,
				TTL:        300,
				Length:     16,
				Data:       &IPv6{netip.MustParseAddr("2607:f8b0:400b:802::2011")},
			},
		},
		{
			name: "Simple A record",
			args: args{
				buf:     []byte("\x06golang\x03com\x00\x00\x01\x00\x01\x00\x00\x01\x2c\x00\x04\x8e\xfb\x29\x51"),
				domains: NewDomains(),
			},
			want: Record{
				Name:   "golang.com",
				Type:   A,
				Class:  1,
				TTL:    300,
				Length: 4,
				Data:   &IPv4{netip.MustParseAddr("142.251.41.81")},
			},
		},
		{
			name: "Bad record type",
			args: args{
				buf:     []byte("\x06golang\x03com\x00\x00\x77\x00\x01\x00\x00\x01\x2c\x00\x04\x8e\xfb\x29\x51"),
				domains: NewDomains(),
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
		TTL        uint32
		Class      uint16
		CacheFlush bool
		Length     uint16
		Type       Type
		Name       string
		Data       RData
	}
	type args struct {
		domains *Domains
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
			args: args{domains: NewDomains()},
			fields: fields{
				TTL:    300,
				Class:  1,
				Length: 16,
				Type:   AAAA,
				Name:   "golang.com",
				Data:   &IPv6{netip.MustParseAddr("2607:f8b0:400b:802::2011")},
			},
			want: []byte("\x06golang\x03com\x00\x00\x1c\x00\x01\x00\x00\x01\x2c\x00\x10\x26\x07\xf8\xb0\x40\x0b\x08\x02\x00\x00\x00\x00\x00\x00\x20\x11"),
		},
		{
			name: "Build Simple AAAA record with mDNS CacheFlush",
			args: args{domains: NewDomains()},
			fields: fields{
				TTL:        300,
				Class:      1,
				CacheFlush: true,
				Length:     16,
				Type:       AAAA,
				Name:       "golang.com",
				Data:       &IPv6{netip.MustParseAddr("2607:f8b0:400b:802::2011")},
			},
			want: []byte("\x06golang\x03com\x00\x00\x1c\x80\x01\x00\x00\x01\x2c\x00\x10\x26\x07\xf8\xb0\x40\x0b\x08\x02\x00\x00\x00\x00\x00\x00\x20\x11"),
		},
		{
			name: "Build Simple A record",
			args: args{domains: NewDomains()},
			fields: fields{
				TTL:    300,
				Class:  1,
				Length: 4,
				Type:   A,
				Name:   "golang.com",
				Data:   &IPv4{netip.MustParseAddr("142.251.41.81")},
			},
			want: []byte("\x06golang\x03com\x00\x00\x01\x00\x01\x00\x00\x01\x2c\x00\x04\x8e\xfb\x29\x51"),
		},
		{
			name: "Build Simple CNAME record",
			args: args{domains: NewDomains()},
			fields: fields{
				TTL:   300,
				Class: 1,
				Type:  CNAME,
				Name:  "golang.com",
				Data:  &CName{Name: "golang.org"},
			},
			want: []byte("\x06golang\x03com\x00\x00\x05\x00\x01\x00\x00\x01\x2c\x00\x0c\x06golang\x03org\x00"),
		},
		{
			name: "Build CNAME record with subdomain pointer",
			args: args{domains: NewDomains()},
			fields: fields{
				TTL:   300,
				Class: 1,
				Type:  CNAME,
				Name:  "golang.com",
				Data:  &CName{Name: "sub.golang.com"},
			},
			want: []byte("\x06golang\x03com\x00\x00\x05\x00\x01\x00\x00\x01\x2c\x00\x06\x03sub\xc0\x00"),
		},
		{
			name: "Build CNAME record with subdomain pointer known issue",
			args: args{domains: NewDomains()},
			fields: fields{
				TTL:   300,
				Class: 1,
				Type:  CNAME,
				Name:  "_acme-challenge.tester.ipv6check.me",
				Data:  &CName{Name: "_acme-challenge.acme.ipv6check.me"},
			},
			want: []byte("\x0f_acme-challenge\x06tester\x09ipv6check\x02me\x00\x00\x05\x00\x01\x00\x00\x01\x2c\x00\x17\x0f_acme-challenge\x04acme\xc0\x17"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Record{
				TTL:        tt.fields.TTL,
				Class:      tt.fields.Class,
				Length:     tt.fields.Length,
				Type:       tt.fields.Type,
				Name:       tt.fields.Name,
				Data:       tt.fields.Data,
				CacheFlush: tt.fields.CacheFlush,
			}
			buf := new(bytes.Buffer)
			if err := r.Build(buf, tt.args.domains); (err != nil) != tt.wantErr {
				t.Errorf("Record.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.want) {
				t.Errorf("Record.Build() = \n%s\n, want \n%s\n", hex.Dump(buf.Bytes()), hex.Dump(tt.want))
			}
		})
	}
}

func BenchmarkRecordParsing(b *testing.B) {
	buf := []byte("\x06golang\x03com\x00\x00\x1c\x00\x01\x00\x00\x01\x2c\x00\x10\x26\x07\xf8\xb0\x40\x0b\x08\x02\x00\x00\x00\x00\x00\x00\x20\x11")
	for i := 0; i < b.N; i++ {
		ParseRecord(bytes.NewBuffer(buf), 0, NewDomains())
	}
}
