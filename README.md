# protoc-gen-psql (PGP)

PGP is a protoc plugin to generate postgesql statement from protobuf files.

This project uses [protoc-gen-star](https://github.com/lyft/protoc-gen-star) to ease code generation.

## How to use it

See `./tests/asset.proto` for example on how to use it.

Two modes are available and may be chosen by adding `alter` as a parameter in protoc command line like `--psql_out="alter:."`

- (default) alter = false : The psql code is generated considering the initial schema is empty (no table, etc.). This mode is great for readability and to generate code for a new database
- alter = true : The psql code is generated to be applied on an existing database. `ALTER` statements are heavily used.

## Tests

Do a `make test-generate` to view a code generated example.

To test the generate psql code, do a `make test`.
