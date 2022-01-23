package dns

import (
	"bytes"
	"testing"
)

func TestMessage_ParseName(t *testing.T) {
	type args struct {
		bytes   []byte
		pointer int
		domains *Domains
	}
	tests := []struct {
		name        string
		args        args
		want        string
		wantErr     bool
		wantDomains map[int]string
	}{
		{
			name: "single domain name",
			args: args{
				bytes:   []byte("\x06domain\x04test\x00"),
				pointer: 0,
				domains: &Domains{parsePtr: map[int]string{}},
			},
			want:        "domain.test",
			wantErr:     false,
			wantDomains: map[int]string{0: "domain.test"},
		},
		{
			name: "domain with sub domain",
			args: args{
				bytes:   []byte("\x03sub\x06domain\x04test\x00"),
				pointer: 0,
				domains: &Domains{parsePtr: map[int]string{}},
			},
			want:        "sub.domain.test",
			wantErr:     false,
			wantDomains: map[int]string{0: "sub.domain.test"},
		},
		{
			name: "domain with sub domain and leading text",
			args: args{
				bytes:   []byte("somerandomtext\x03sub\x06domain\x04test\x00"),
				pointer: 14,
				domains: &Domains{parsePtr: map[int]string{}},
			},
			want:        "sub.domain.test",
			wantErr:     false,
			wantDomains: map[int]string{14: "sub.domain.test"},
		},
		{
			name: "domain with pointer",
			args: args{
				bytes:   []byte("somerandomtext\x03sub\x06domain\x04test\x00\xc0\x0e"),
				pointer: 31,
				domains: &Domains{parsePtr: map[int]string{14: "sub.domain.test"}},
			},
			want:        "sub.domain.test",
			wantErr:     false,
			wantDomains: map[int]string{14: "sub.domain.test"},
		},
		{
			name: "domain bad with pointer",
			args: args{
				bytes:   []byte("somerandomtext\x03sub\x06domain\x04test\x00\xc0\x0f"),
				pointer: 31,
				domains: &Domains{parsePtr: map[int]string{14: "sub.domain.test"}},
			},
			want:        "",
			wantErr:     true,
			wantDomains: map[int]string{14: "sub.domain.test"},
		},
		{
			name: "bad pointer to domain",
			args: args{
				bytes:   []byte("somerandomtext\x03sub\x06domain\x04test\x00"),
				pointer: 15,
				domains: &Domains{parsePtr: map[int]string{}},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "bad length indicator",
			args: args{
				bytes:   []byte("somerandomtext\x03sub\xff\x04test\x00"),
				pointer: 14,
				domains: &Domains{parsePtr: map[int]string{}},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "root domain",
			args: args{
				bytes:   []byte("\x00"),
				domains: NewDomains(),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseName(bytes.NewBuffer(tt.args.bytes[tt.args.pointer:]),
				tt.args.pointer, tt.args.domains)
			if (err != nil) != tt.wantErr {
				t.Errorf("Message.ParseName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Message.ParseName() = %q, want %q", got, tt.want)
				return
			}
			for k, v := range tt.wantDomains {
				if vv, ok := tt.args.domains.GetParse(k); !ok || vv != v {
					t.Errorf("Message.ParseName() domains[%d] = %v, want %v", k, v, vv)
				}
			}
		})
	}
}

func TestBuildName(t *testing.T) {
	type args struct {
		name    string
		domains *Domains
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "single domain name",
			args: args{
				name:    "domain.test",
				domains: &Domains{buildPtr: map[string]int{}},
			},
			want: "\x06domain\x04test\x00",
		},
		{
			name: "sub domain name",
			args: args{
				name:    "sub.domain.test",
				domains: &Domains{buildPtr: map[string]int{}},
			},
			want: "\x03sub\x06domain\x04test\x00",
		},
		{
			name: "cached domain name",
			args: args{
				name:    "domain.test",
				domains: &Domains{buildPtr: map[string]int{"domain.test": 42}},
			},
			want: "\xc0\x2a",
		},
		{
			name: "mixed domain name",
			args: args{
				name:    "sub.domain.test",
				domains: &Domains{buildPtr: map[string]int{"domain.test": 42}},
			},
			want: "\x03sub\xc0\x2a",
		},
		{
			name: "root domain",
			args: args{
				name:    "",
				domains: NewDomains(),
			},
			want: "\x00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildName(tt.args.name, tt.args.domains)
			if got != tt.want {
				t.Errorf("BuildName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
