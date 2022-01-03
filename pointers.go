package dnsmessage

type Pointers struct {
	parsePtr map[int]string
	buildPtr map[string]int
}

func (p *Pointers) GetParse(ptr int) (string, bool) {
	name, ok := p.parsePtr[ptr]
	return name, ok
}

func (p *Pointers) GetBuild(name string) (int, bool) {
	ptr, ok := p.buildPtr[name]
	return ptr, ok
}

func (p *Pointers) SetParse(ptr int, name string) {
	ok := true
	for i, c := range name {
		if ok {
			p.parsePtr[ptr+i] = name[i:len(name)]
			ok = false
		} else if c == '.' {
			ok = true
		}
	}
}

func (p *Pointers) SetBuild(ptr int, name string) {
	ok := true
	for i, c := range name {
		if ok {
			p.buildPtr[name[i:len(name)]] = ptr + i
			ok = false
		} else if c == '.' {
			ok = true
		}
	}
}

func NewPointers() *Pointers {
	return &Pointers{parsePtr: map[int]string{}, buildPtr: map[string]int{}}
}
