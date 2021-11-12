package dnsmessage

import (
	"testing"
)

func TestMessage_ParseName(t *testing.T) {
	type args struct {
		bytes   []byte
		pointer int
		domains map[int]string
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
				domains: map[int]string{},
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
				domains: map[int]string{},
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
				domains: map[int]string{},
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
				domains: map[int]string{14: "sub.domain.test"},
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
				domains: map[int]string{14: "sub.domain.test"},
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
				domains: map[int]string{},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "bad pointer argument",
			args: args{
				bytes:   []byte("somerandomtext\x03sub\x06domain\x04test\x00"),
				pointer: 200,
				domains: map[int]string{},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "bad length indicator",
			args: args{
				bytes:   []byte("somerandomtext\x03sub\xff\x04test\x00"),
				pointer: 14,
				domains: map[int]string{},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseName(tt.args.bytes, tt.args.pointer, tt.args.domains)
			if (err != nil) != tt.wantErr {
				t.Errorf("Message.ParseName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Message.ParseName() = %q, want %q", got, tt.want)
				return
			}
			for k, v := range tt.wantDomains {
				if vv, ok := tt.args.domains[k]; !ok || vv != v {
					t.Errorf("Message.ParseName() pointers[%d] = %v, want %v", k, v, vv)
				}
			}
		})
	}
}
