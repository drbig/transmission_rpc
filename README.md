# transmission_rpc [![GoDoc](https://godoc.org/github.com/drbig/transmission_rpc?status.svg)](http://godoc.org/github.com/drbig/transmission_rpc)

Package transmission_rpc provides unassuming interface to the Transmission RPC service.

Features / bugs:

- Concurrent-usage safe
- Provides raw `[]byte`-based access and minimal semantic interface
- Handles optional authentication
- Semantic requests are `tag` and `success` guarded
- Self-contained package, no external dependencies
- No test suite (for now)
- No example cmd (for now)

## Contributing

Follow the usual GitHub development model:

1. Clone the repository
2. Make your changes on a separate branch
3. Make sure you run `gofmt` and `go test` before committing
4. Make a pull request

See licensing for legalese.

## Licensing

Standard two-clause BSD license, see LICENSE.txt for details.

Any contributions will be licensed under the same conditions.

Copyright (c) 2014 Piotr S. Staszewski
