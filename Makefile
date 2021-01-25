empty :=
space := $(empty) $(empty)
NAME := psql
PACKAGE := github.com/intrinsec/protoc-gen-$(NAME)

CI_JOB_ID ?= local
export CI_JOB_ID

# protoc-gen-go parameters for properly generating the import path for PGV
GO_IMPORT_SPACES := M$(NAME)/$(NAME).proto=${PACKAGE}/$(NAME),\
	Mgoogle/protobuf/any.proto=github.com/golang/protobuf/ptypes/any,\
	Mgoogle/protobuf/duration.proto=github.com/golang/protobuf/ptypes/duration,\
	Mgoogle/protobuf/struct.proto=github.com/golang/protobuf/ptypes/struct,\
	Mgoogle/protobuf/timestamp.proto=github.com/golang/protobuf/ptypes/timestamp,\
	Mgoogle/protobuf/wrappers.proto=github.com/golang/protobuf/ptypes/wrappers,\
	Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor
GO_IMPORT:=$(subst $(space),,$(GO_IMPORT_SPACES))

.PHONY: build
build: bin/protoc-gen-$(NAME)

.PHONY: install
install: $(NAME)/$(NAME).pb.go
	@go install -v .

$(NAME)/$(NAME).pb.go: bin/protoc-gen-go $(NAME)/$(NAME).proto
	@cd $(NAME) && protoc -I . \
		--plugin=protoc-gen-go=$(shell pwd)/bin/protoc-gen-go \
		--go_opt=paths=source_relative \
		--go_out="${GO_IMPORT}:." $(NAME).proto

bin/protoc-gen-go:
	@GOBIN=$(shell pwd)/bin go install google.golang.org/protobuf/cmd/protoc-gen-go


bin/protoc-gen-$(NAME): $(NAME)/$(NAME).pb.go $(wildcard *.go)
	@GOBIN=$(shell pwd)/bin go install .


.PHONY: test-generate
test-generate: build
	@protoc -I . --plugin=protoc-gen-$(NAME)=$(shell pwd)/bin/protoc-gen-$(NAME) --$(NAME)_out="." tests/asset.proto
	@cat tests/*.pb.psql


.PHONY: build-docker
build-docker:
	docker-compose -p $(NAME)-$(CI_JOB_ID) -f ./tests/docker-compose.tests.yml build

.PHONY: test-integration
test-integration: build
	docker-compose -p $(NAME)-$(CI_JOB_ID) -f ./tests/docker-compose.tests.yml rm -f
	docker-compose -p $(NAME)-$(CI_JOB_ID) -f ./tests/docker-compose.tests.yml up --exit-code-from=client
	docker-compose -p $(NAME)-$(CI_JOB_ID) -f ./tests/docker-compose.tests.yml down -v


.PHONY: test
test: build test-generate test-integration


.PHONY: clean
clean:
	@rm -fv tests/*.psql


.PHONY: distclean
distclean: clean
	@rm -fv bin/protoc-gen-go bin/protoc-gen-$(NAME) $(NAME)/$(NAME).pb.go

