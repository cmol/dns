package dns

type Type uint16
type Class uint16

const (
	OptRa = 0x80
	OptRd = 0x100
	OptTc = 0x200
	OptAa = 0x400
	OptQr = 0x8000
)

const (
	IN Class = 1
)

const (
	A          Type = 1
	NS         Type = 2
	CNAME      Type = 5
	SOA        Type = 6
	PTR        Type = 12
	HINFO      Type = 13
	MX         Type = 15
	TXT        Type = 16
	RP         Type = 17
	AFSDB      Type = 18
	SIG        Type = 24
	KEY        Type = 25
	AAAA       Type = 28
	LOC        Type = 29
	SRV        Type = 33
	NAPTR      Type = 35
	KX         Type = 36
	CERT       Type = 37
	DNAME      Type = 39
	OPT        Type = 41
	APL        Type = 42
	DS         Type = 43
	SSHFP      Type = 44
	IPSECKEY   Type = 45
	RRSIG      Type = 46
	NSE        Type = 47
	DNSKEY     Type = 48
	DHCID      Type = 49
	NSEC3      Type = 50
	NSEC3PARAM Type = 51
	TLSA       Type = 52
	SMIMEA     Type = 53
	HIP        Type = 55
	CDS        Type = 59
	CDNSKEY    Type = 60
	OPENPGPKEY Type = 61
	CSYNC      Type = 62
	ZONEMD     Type = 63
	SVCB       Type = 64
	HTTPS      Type = 65
	EUI48      Type = 108
	EUI64      Type = 109
	TKEY       Type = 249
	TSIG       Type = 250
	IXFR       Type = 251
	AXFR       Type = 252
	URI        Type = 256
	CAA        Type = 257
	TA         Type = 32768
	DLV        Type = 32769
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
