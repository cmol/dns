package dnsmessage

type RData interface {
	Parse()
	Build()
	RecordType() uint16
}

type RR struct {
	ID, QDCount, ANCount, NSCount, ARCount uint16
	OpCode, Opts, RCode                    uint8
	QR                                     bool
	Data                                   *RData
}
