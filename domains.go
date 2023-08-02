package dns

// Domains holds maps for parsing and building CNAME pointers
type Domains struct {
	parsePtr map[int]string
	buildPtr map[string]int
}

// GetParse returns a domain for a given pointer
func (p *Domains) GetParse(ptr int) (string, bool) {
	name, ok := p.parsePtr[ptr]
	return name, ok
}

// GetBuild returns af pointer for a given name
func (p *Domains) GetBuild(name string) (int, bool) {
	ptr, ok := p.buildPtr[name]
	return ptr, ok
}

// SetParse adds parse pointers to the domain map
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

// NewDomains returns a new map for parsing and building
func NewDomains() *Domains {
	return &Domains{parsePtr: map[int]string{}, buildPtr: map[string]int{}}
}
