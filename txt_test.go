package dns

import (
	"bytes"
	"reflect"
	"testing"
)

func TestTxt_Parse(t *testing.T) {
	type args struct {
		buf        []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Txt
		wantErr bool
	}{
		{
			name: "Simple domain name",
			args: args{
				buf:        []byte("\x45v=spf1 ip4:192.0.2.0/24 ip4:198.51.100.123 ip6:2620:0:860::/46 a -all"),
			},
			want: Txt{
        Length: 69,
        Data: "v=spf1 ip4:192.0.2.0/24 ip4:198.51.100.123 ip6:2620:0:860::/46 a -all",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
      domains := NewDomains()
			n := &Txt{}
			if err := n.Parse(bytes.NewBuffer(tt.args.buf), 0, domains); (err != nil) != tt.wantErr {
				t.Errorf("Txt.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(n, &tt.want) {
				t.Errorf("Txt.Parse() = \n%+v\n, want \n%+v", n, tt.want)
			}
		})
	}
}

func TestTxt_Build(t *testing.T) {
	tests := []struct {
		name       string
		txt        Txt
		want       []byte
		wantLength int
		wantErr    bool
	}{
		{
			name: "Simple TXT SFP record",
			txt: Txt{
        Length: 69,
        Data: "v=spf1 ip4:192.0.2.0/24 ip4:198.51.100.123 ip6:2620:0:860::/46 a -all",
			},
			want:       []byte("\x45v=spf1 ip4:192.0.2.0/24 ip4:198.51.100.123 ip6:2620:0:860::/46 a -all"),
			wantLength: 70,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			if length, err := tt.txt.PreBuild(&Record{}, NewDomains()); err != nil || length != tt.wantLength {
				if err != nil {
					t.Errorf("Txt.PreBuild() error = %v", err)
				} else {
					t.Errorf("Txt.PreBuild() = %v, want %v", length, tt.wantLength)
				}
				return
			}
			if err := tt.txt.Build(buf, NewDomains()); (err != nil) != tt.wantErr {
				t.Errorf("Txt.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(buf.Bytes(), tt.want) {
				t.Errorf("Txt.Build() = \n%+v, want \n%+v", buf.Bytes(), tt.want)
			}
		})
	}
}
