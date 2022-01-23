package dns

import (
	"bytes"
	"reflect"
	"testing"
)

func TestOpt_Parse(t *testing.T) {
	tests := []struct {
		name    string
		record  *Record
		buf     []byte
		want    *Opt
		wantErr bool
	}{
		{
			name: "Simple OPT record",
			record: &Record{
				Class: 512,
				TTL:   0,
			},
			buf: []byte{},
			want: &Opt{
				UDPSize: 512,
				Options: map[uint16][]byte{},
			},
		},
		{
			name: "Advanced OPT record",
			record: &Record{
				Class: 512,
				TTL:   0xb2038000,
			},
			buf: []byte{},
			want: &Opt{
				UDPSize:     512,
				RCode:       0xb2,
				EDNSVersion: 0x03,
				DNSSec:      true,
				Options:     map[uint16][]byte{},
			},
		},
		{
			name: "Opt extra options",
			record: &Record{
				Class:  512,
				Length: 6,
			},
			buf: []byte("\x00\x05\x00\x02\xab\xab"),
			want: &Opt{
				UDPSize: 512,
				Options: map[uint16][]byte{5: []byte("\xab\xab")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Opt{Record: tt.record}
			if err := o.Parse(bytes.NewBuffer(tt.buf), 0, NewDomains()); (err != nil) != tt.wantErr {
				t.Errorf("Opt.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.want.Record = tt.record
			if !reflect.DeepEqual(o, tt.want) {
				t.Errorf("Opt.Parse() = %v, want %v", o, tt.want)
			}
		})
	}
}

func TestOpt_PreBuild(t *testing.T) {
	type fields struct {
		UDPSize     uint16
		RCode       byte
		EDNSVersion byte
		DNSSec      bool
		Record      *Record
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Record
		wantErr bool
	}{
		{
			name: "OPT record",
			fields: fields{
				UDPSize:     512,
				RCode:       0xb2,
				EDNSVersion: 0x03,
				DNSSec:      true,
			},
			want: &Record{
				Class: 512,
				TTL:   0xb2038000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Opt{
				UDPSize:     tt.fields.UDPSize,
				RCode:       tt.fields.RCode,
				EDNSVersion: tt.fields.EDNSVersion,
				DNSSec:      tt.fields.DNSSec,
			}
			r := &Record{}
			if _, err := o.PreBuild(r, NewDomains()); (err != nil) != tt.wantErr {
				t.Errorf("Opt.PreBuild() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(r, tt.want) {
				t.Errorf("Opt.PreBuild() = %v, want %v", r, tt.want)
			}
		})
	}
}

func TestDefaultOpt(t *testing.T) {
	tests := []struct {
		name string
		size int
		want []byte
	}{
		{
			name: "Simple OPT",
			size: 1680,
			want: []byte("\x00\x00\x29\x06\x90\x00\x00\x00\x00\x00\x00"),
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			opt := DefaultOpt(tt.size)
			if opt.Build(buf, NewDomains()); !reflect.DeepEqual(buf.Bytes(), tt.want) {
				t.Errorf("DefaultOpt() = %v, want %v", buf.Bytes(), tt.want)
			}
		})
	}
}
