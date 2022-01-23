# DNS

DNS is a hands on library for DNS message parsing and building in Golang for
those who would like to control the details of DNS messages.

## Prerequisites

This library is not a DNS server. You can build a DNS server with this library,
but it is not handled by the library itself. This means that you need a basic
knowledge of DNS to use this library, though examples will be provided.

## Installation

Add import to your application:

```
import "github.com/cmol/dns"
```

## Usage

```golang
// Read a DNS message query from a 'bytes.Buffer'
query, err = dns.ParseMessage(buffer)

// Make a reply struct to said message
reply = dns.ReplyTo(query)

// Add EDNS information to the message
reply.Additional = []dns.Record{*dns.DefaultOpt(1024)}

// Create output buffer and build message in the buffer
buf := new(bytes.Buffer)
err = reply.Build(buf, dns.NewDomains())

// Send out the message on an existing remote connection
n, err = connection.WriteToUDP(buf.Bytes(), remoteAddr)
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to
discuss what you would like to change.

New Resource Record parsers and builders are welcome, but 

Please make sure to update tests as appropriate.


## Supported Resource Record types

Current list of supported RR types, though I will add more as I use the
library:

 - A
 - AAAA
 - CNAME
 - OPT

More types are welcome! Please add tests and references to the standard in the
PR.

## License
[MIT](https://choosealicense.com/licenses/mit/)
