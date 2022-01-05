package dnsmessage

import (
	"testing"
)

func TestDomains_SetParse(t *testing.T) {
	type fields struct {
		parsePtr map[int]string
	}
	type args struct {
		ptr  int
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[int]string
	}{
		{
			name:   "Simple domain",
			fields: fields{parsePtr: map[int]string{}},
			args: args{
				ptr:  0,
				name: "domain.test",
			},
			want: map[int]string{0: "domain.test", 7: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Domains{
				parsePtr: tt.fields.parsePtr,
			}
			d.SetParse(tt.args.ptr, tt.args.name)
			for k, v := range tt.want {
				if vv, ok := d.GetParse(k); !ok || vv != v {
					t.Errorf("Domains.SetParse() domains[%v] = %v, want %v", k, v, vv)
				}
			}
		})
	}
}

func TestDomains_SetBuild(t *testing.T) {
	type fields struct {
		buildPtr map[string]int
	}
	type args struct {
		ptr  int
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]int
	}{
		{
			name:   "Simple domain",
			fields: fields{buildPtr: map[string]int{}},
			args: args{
				ptr:  0,
				name: "domain.test",
			},
			want: map[string]int{"domain.test": 0, "test": 7},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Domains{
				buildPtr: tt.fields.buildPtr,
			}
			d.SetBuild(tt.args.ptr, tt.args.name)
			for k, v := range tt.want {
				if vv, ok := d.GetBuild(k); !ok || vv != v {
					t.Errorf("Domains.SetBuild() domains[%v] = %v, want %v", k, v, vv)
				}
			}
		})
	}
}
