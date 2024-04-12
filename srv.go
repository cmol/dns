package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

// Srv implements interface RData
type Srv struct {
	Priority    uint16
	Weight      uint16
	Port        uint16
	Target      string
	targetBytes string
	NameBytes   string
	// It seems like a lot of services does things like
	//   [identifier]._[service]._[proto].[name]
	// Support this scheme by adding data to this field if it there's more parts
	// to be parsed from the name
	Identifier string
	Service    string
	Proto      string
	Name       string
}

// Parse implements A parsing for interface RData
func (s *Srv) Parse(buf *bytes.Buffer, ptr int, domains *Domains) error {
	if err := s.parseName(); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &s.Priority); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &s.Weight); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &s.Port); err != nil {
		return err
	}
	name, err := ParseName(buf, ptr+6, domains)
	if err != nil {
		return errors.New("unable to parse Srv: " + err.Error())
	}
	s.Target = name
	return nil
}

// Build implements A building for interface RData
func (s *Srv) Build(buf *bytes.Buffer, domains *Domains) error {
	if err := binary.Write(buf, binary.BigEndian, s.Priority); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, s.Weight); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, s.Port); err != nil {
		return err
	}
	domains.SetBuild(buf.Len(), s.Target)
	buf.WriteString(s.targetBytes)
	return nil
}

// PreBuild step, building name and adding full record
func (s *Srv) PreBuild(_ *Record, domains *Domains) (int, error) {
	s.targetBytes = BuildName(s.Target, domains)
	return len(s.targetBytes) + 6, nil
}

// TransformName adds service/proto/name fields of server record
func (s *Srv) TransformName(name string) string {
	generatedName := ""
	if len(s.Identifier) > 1 {
		generatedName = s.Identifier + "."
	}
	s.NameBytes = generatedName + s.Service + "." + s.Proto + "." + s.Name
	return s.NameBytes
}

func (s *Srv) parseName() error {
	parts := strings.Split(s.NameBytes, ".")
	pLen := len(parts)
	if pLen < 3 {
		return errors.New("not enough name parts of SRV record in: " + s.NameBytes)
	}
	// mDNS records seems to have this extra field for DNS-SD (with name almost
	// always being `local`). Put it into a specific field instead of folding it
	// into the service.
	if pLen > 3 {
		s.Identifier = strings.Join(parts[0:pLen-3], ".")
	}
	s.Service = parts[pLen-3]
	s.Proto = parts[pLen-2]
	s.Name = parts[pLen-1]

	return nil
}
