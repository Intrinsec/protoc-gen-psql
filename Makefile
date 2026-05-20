empty :=
space := $(empty) $(empty)
NAME := psql
PACKAGE := github.com/intrinsec/protoc-gen-$(NAME)

SHELL := /bin/bash

CI_JOB_ID ?= local
export CI_JOB_ID

# protoc-gen-go import remap for the local psql package.
# Well-known types resolve via google.golang.org/protobuf/types/* by default —
# no M-mappings required.
GO_IMPORT_SPACES := M$(NAME)/$(NAME).proto=${PACKAGE}/$(NAME)
GO_IMPORT:=$(subst $(space),,$(GO_IMPORT_SPACES))

.PHONY: build
build: bin/protoc-gen-$(NAME)

.PHONY: install
install: $(NAME)/$(NAME).pb.go
	@go install -mod=vendor -v .

# Resolve protoc's well-known type include dir (e.g. google/protobuf/descriptor.proto).
# `protobuf-compiler` on Debian/Ubuntu installs them under /usr/include; on macOS
# (Homebrew) they live under $(brew --prefix)/include. Fall back to /usr/include.
PROTOC_WKT_INCLUDE ?= $(shell test -d /opt/homebrew/include/google/protobuf && echo /opt/homebrew/include || echo /usr/include)

$(NAME)/$(NAME).pb.go: bin/protoc-gen-go $(NAME)/$(NAME).proto
	@cd $(NAME) && protoc -I . -I $(PROTOC_WKT_INCLUDE) \
		--plugin=protoc-gen-go=$(shell pwd)/bin/protoc-gen-go \
		--go_opt=paths=source_relative \
		--go_out="${GO_IMPORT}:." $(NAME).proto

bin/protoc-gen-go:
	@GOBIN=$(shell pwd)/bin go install -mod=mod google.golang.org/protobuf/cmd/protoc-gen-go


bin/protoc-gen-$(NAME): $(NAME)/$(NAME).pb.go $(wildcard *.go)
	@GOBIN=$(shell pwd)/bin go install -mod=vendor .

.PHONY: vendor
vendor:
	@go mod tidy
	@go mod vendor


.PHONY: test-generate
test-generate: build
	@protoc -I . -I $(PROTOC_WKT_INCLUDE) --plugin=protoc-gen-$(NAME)=$(shell pwd)/bin/protoc-gen-$(NAME) --$(NAME)_out="." tests/*.proto
	@cat tests/*.pb.psql
	# Checking diff between generated file and reference file
	# If the following is empty then the file are identical
	@for i in `ls tests/*.pb.psql`; do \
		diff tests/references/$$(basename $$i) $$i; \
	done \


# `docker-compose` (v1, hyphenated) is reaching EOL; prefer `docker compose`
# (v2 plugin). Fall back to the legacy binary if v2 is not installed.
DOCKER_COMPOSE ?= $(shell docker compose version >/dev/null 2>&1 && echo 'docker compose' || echo 'docker-compose')

.PHONY: build-docker
build-docker:
	$(DOCKER_COMPOSE) -p $(NAME)-$(CI_JOB_ID) -f ./tests/docker-compose.tests.yml build

.PHONY: test-integration
test-integration: build
	$(DOCKER_COMPOSE) -p $(NAME)-$(CI_JOB_ID) -f ./tests/docker-compose.tests.yml rm -f
	$(DOCKER_COMPOSE) -p $(NAME)-$(CI_JOB_ID) -f ./tests/docker-compose.tests.yml up --exit-code-from=client
	$(DOCKER_COMPOSE) -p $(NAME)-$(CI_JOB_ID) -f ./tests/docker-compose.tests.yml down -v


.PHONY: test
test: build test-generate test-integration


.PHONY: clean
clean:
	@rm -fv tests/*.psql


.PHONY: distclean
distclean: clean
	@rm -fv bin/protoc-gen-go bin/protoc-gen-$(NAME) $(NAME)/$(NAME).pb.go
