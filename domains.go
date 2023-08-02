package dns

type Domains struct {
	parsePtr map[int]string
	buildPtr map[string]int
}

func (p *Domains) GetParse(ptr int) (string, bool) {
	name, ok := p.parsePtr[ptr]
	return name, ok
}

func (p *Domains) GetBuild(name string) (int, bool) {
	ptr, ok := p.buildPtr[name]
	return ptr, ok
}

func (p *Domains) SetParse(ptr int, name string) {
	ok := true
	for i, c := range name {
		if ok {
			p.parsePtr[ptr+i] = name[i:]
			ok = false
		} else if c == '.' {
			ok = true
		}
	}
}

// SetBuild adds build pointers to the domain map
func (p *Domains) SetBuild(ptr int, name string) {
	ok := true
	for i, c := range name {
		if ok {
			if _, found := p.buildPtr[name[i:]]; !found {
				p.buildPtr[name[i:]] = ptr + i
			}
			ok = false
		} else if c == '.' {
			ok = true
		}
	}
}

func NewDomains() *Domains {
	return &Domains{parsePtr: map[int]string{}, buildPtr: map[string]int{}}
}
