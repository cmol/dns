package dns

import (
	"bytes"
	"reflect"
	"testing"
)

func TestSrv_Parse(t *testing.T) {
	type args struct {
		buf        []byte
		ptr        int
		domains    *Domains
		recordName string
	}
	tests := []struct {
		name    string
		args    args
		want    Srv
		wantErr bool
	}{
		{
			name: "Simple domain name",
			args: args{
				buf:        []byte("\x00\x05\x00\x0a\x1b\x58\x0c\x58\x30\x30\x31\x30\x30\x4e\x37\x47\x31\x45\x59\x05local\x00"),
				domains:    NewDomains(),
				recordName: "50in Hisense Roku TV._airplay._tcp.local",
			},
			want: Srv{
				Priority:   5,
				Weight:     10,
				Port:       7000,
				Target:     "X00100N7G1EY.local",
				Identifier: "50in Hisense Roku TV",
				Service:    "_airplay",
				Proto:      "_tcp",
				Name:       "local",
				NameBytes:  "50in Hisense Roku TV._airplay._tcp.local",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Srv{NameBytes: tt.args.recordName}
			if err := n.Parse(bytes.NewBuffer(tt.args.buf), tt.args.ptr, tt.args.domains); (err != nil) != tt.wantErr {
				t.Errorf("Srv.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(n, &tt.want) {
				t.Errorf("Srv.Parse() = \n%+v\n, want \n%+v", n, tt.want)
			}
		})
	}
}

func TestSrv_Build(t *testing.T) {
	tests := []struct {
		name       string
		srv        Srv
		want       []byte
		wantLength int
		wantErr    bool
	}{
		{
			name: "Simple domain",
			srv: Srv{
				Priority:   5,
				Weight:     10,
				Port:       7000,
				Target:     "X00100N7G1EY.local",
				Identifier: "50in Hisense Roku TV",
				Service:    "_airplay",
				Proto:      "_tcp",
				Name:       "local",
				NameBytes:  "50in Hisense Roku TV._airplay._tcp.local",
			},
			want:       []byte("\x00\x05\x00\x0a\x1b\x58\x0c\x58\x30\x30\x31\x30\x30\x4e\x37\x47\x31\x45\x59\x05local\x00"),
			wantLength: 26,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			if length, err := tt.srv.PreBuild(&Record{}, NewDomains()); err != nil || length != tt.wantLength {
				if err != nil {
					t.Errorf("Srv.PreBuild() error = %v", err)
				} else {
					t.Errorf("Srv.PreBuild() = %v, want %v", length, tt.wantLength)
				}
				return
			}
			if err := tt.srv.Build(buf, NewDomains()); (err != nil) != tt.wantErr {
				t.Errorf("Srv.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.want) {
				t.Errorf("Srv.Build() = \n%+v, want \n%+v", buf.Bytes(), tt.want)
			}
		})
	}
}

func TestSrv_TransformName(t *testing.T) {
	tests := []struct {
		name     string
		srv      Srv
		wantName string
	}{
		{
			name: "Simple transform",
			srv: Srv{
				Priority:   5,
				Weight:     10,
				Port:       7000,
				Target:     "X00100N7G1EY.local",
				Identifier: "50in Hisense Roku TV",
				Service:    "_airplay",
				Proto:      "_tcp",
				Name:       "local",
				NameBytes:  "50in Hisense Roku TV._airplay._tcp.local",
			},
			wantName: "50in Hisense Roku TV._airplay._tcp.local",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := tt.srv.TransformName("")
			if name != tt.wantName {
				t.Errorf("Srv.TransformName() = %v, want %v", name, tt.wantName)
			}
		})
	}
}
