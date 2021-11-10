package dnsmessage

import (
	"testing"
)

func TestMessage_ParseName(t *testing.T) {
	type fields struct {
		pointers map[int]string
	}
	type args struct {
		bytes   []byte
		pointer int
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		want         string
		wantErr      bool
		wantPointers map[int]string
	}{
		{
			name:   "single domain name",
			fields: fields{pointers: map[int]string{}},
			args: args{
				bytes:   []byte("\x06domain\x04test\x00"),
				pointer: 0,
			},
			want:         "domain.test",
			wantErr:      false,
			wantPointers: map[int]string{0: "domain.test"},
		},
		{
			name:   "domain with sub domain",
			fields: fields{pointers: map[int]string{}},
			args: args{
				bytes:   []byte("\x03sub\x06domain\x04test\x00"),
				pointer: 0,
			},
			want:         "sub.domain.test",
			wantErr:      false,
			wantPointers: map[int]string{0: "sub.domain.test"},
		},
		{
			name:   "domain with sub domain and leading text",
			fields: fields{pointers: map[int]string{}},
			args: args{
				bytes:   []byte("somerandomtext\x03sub\x06domain\x04test\x00"),
				pointer: 14,
			},
			want:         "sub.domain.test",
			wantErr:      false,
			wantPointers: map[int]string{14: "sub.domain.test"},
		},
		{
			name:   "bad pointer to domain",
			fields: fields{pointers: map[int]string{}},
			args: args{
				bytes:   []byte("somerandomtext\x03sub\x06domain\x04test\x00"),
				pointer: 15,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				pointers: tt.fields.pointers,
			}
			got, err := m.ParseName(tt.args.bytes, tt.args.pointer)
			if (err != nil) != tt.wantErr {
				t.Errorf("Message.ParseName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Message.ParseName() = %q, want %q", got, tt.want)
				return
			}
			for k, v := range tt.wantPointers {
				if vv, ok := m.pointers[k]; !ok || vv != v {
					t.Errorf("Message.ParseName() pointers[%d] = %v, want %v", k, v, vv)
				}
			}
		})
	}
}
