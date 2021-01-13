empty :=
space := $(empty) $(empty)
PACKAGE := github.com/intrinsec/protoc-gen-psql

# protoc-gen-go parameters for properly generating the import path for PGV
PSQL_IMPORT := Mpsql/psql.proto=${PACKAGE}/psql
GO_IMPORT_SPACES := ${PSQL_IMPORT},\
	Mgoogle/protobuf/any.proto=github.com/golang/protobuf/ptypes/any,\
	Mgoogle/protobuf/duration.proto=github.com/golang/protobuf/ptypes/duration,\
	Mgoogle/protobuf/struct.proto=github.com/golang/protobuf/ptypes/struct,\
	Mgoogle/protobuf/timestamp.proto=github.com/golang/protobuf/ptypes/timestamp,\
	Mgoogle/protobuf/wrappers.proto=github.com/golang/protobuf/ptypes/wrappers,\
	Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor
GO_IMPORT:=$(subst $(space),,$(GO_IMPORT_SPACES))

.PHONY: build
build: psql/psql.pb.go
	@GOBIN=$(shell pwd)/bin go install -v .

psql/psql.pb.go: bin/protoc-gen-go psql/psql.proto
	@cd psql && protoc -I . \
		--plugin=protoc-gen-go=$(shell pwd)/bin/protoc-gen-go \
		--go_opt=paths=source_relative \
		--go_out="${GO_IMPORT}:." psql.proto

bin/protoc-gen-go:
	@GOBIN=$(shell pwd)/bin go install github.com/golang/protobuf/protoc-gen-go


.PHONY: test
test: build
	@PATH="$$(pwd)/bin:$(PATH)" protoc -I . --psql_out="." asset.proto
	@cat asset.psql
