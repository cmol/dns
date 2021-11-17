package dnsmessage

type Type uint16
type Class uint16

const (
	IN Class = 1
)

const (
	A          Type = 1
	NS              = 2
	CNAME           = 5
	SOA             = 6
	PTR             = 12
	HINFO           = 13
	MX              = 15
	TXT             = 16
	RP              = 17
	AFSDB           = 18
	SIG             = 24
	KEY             = 25
	AAAA            = 28
	LOC             = 29
	SRV             = 33
	NAPTR           = 35
	KX              = 36
	CERT            = 37
	DNAME           = 39
	OPT             = 41
	APL             = 42
	DS              = 43
	SSHFP           = 44
	IPSECKEY        = 45
	RRSIG           = 46
	NSE             = 47
	DNSKEY          = 48
	DHCID           = 49
	NSEC3           = 50
	NSEC3PARAM      = 51
	TLSA            = 52
	SMIMEA          = 53
	HIP             = 55
	CDS             = 59
	CDNSKEY         = 60
	OPENPGPKEY      = 61
	CSYNC           = 62
	ZONEMD          = 63
	SVCB            = 64
	HTTPS           = 65
	EUI48           = 108
	EUI64           = 109
	TKEY            = 249
	TSIG            = 250
	IXFR            = 251
	AXFR            = 252
	URI             = 256
	CAA             = 257
	TA              = 32768
	DLV             = 32769
)

var RRTypeStrings = map[Type]string{
	A:          "A",
	NS:         "NS",
	CNAME:      "CNAME",
	SOA:        "SOA",
	PTR:        "PTR",
	HINFO:      "HINFO",
	MX:         "MX",
	TXT:        "TXT",
	RP:         "RP",
	AFSDB:      "AFSDB",
	SIG:        "SIG",
	KEY:        "KEY",
	AAAA:       "AAAA",
	LOC:        "LOC",
	SRV:        "SRV",
	NAPTR:      "NAPTR",
	KX:         "KX",
	CERT:       "CERT",
	DNAME:      "DNAME",
	OPT:        "OPT",
	APL:        "APL",
	DS:         "DS",
	SSHFP:      "SSHFP",
	IPSECKEY:   "IPSECKEY",
	RRSIG:      "RRSIG",
	NSE:        "NSE",
	DNSKEY:     "DNSKEY",
	DHCID:      "DHCID",
	NSEC3:      "NSEC3",
	NSEC3PARAM: "NSEC3PARAM",
	TLSA:       "TLSA",
	SMIMEA:     "SMIMEA",
	HIP:        "HIP",
	CDS:        "CDS",
	CDNSKEY:    "CDNSKEY",
	OPENPGPKEY: "OPENPGPKEY",
	CSYNC:      "CSYNC",
	ZONEMD:     "ZONEMD",
	SVCB:       "SVCB",
	HTTPS:      "HTTPS",
	EUI48:      "EUI48",
	EUI64:      "EUI64",
	TKEY:       "TKEY",
	TSIG:       "TSIG",
	IXFR:       "IXFR",
	AXFR:       "AXFR",
	URI:        "URI",
	CAA:        "CAA",
	TA:         "TA",
	DLV:        "DLV",
}
