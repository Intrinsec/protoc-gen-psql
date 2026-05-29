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


PROTO_FIXTURES := $(sort $(shell find tests -name '*.proto' -not -path 'tests/buf/*'))

.PHONY: test-generate
test-generate: build
	@protoc -I . -I $(PROTOC_WKT_INCLUDE) --plugin=protoc-gen-$(NAME)=$(shell pwd)/bin/protoc-gen-$(NAME) --$(NAME)_out="." $(PROTO_FIXTURES)
	@find tests -name '*.pb.psql' -not -path 'tests/references/*' -exec cat {} \;
	# Checking diff between generated file and reference file
	# If the following is empty then the file are identical
	@for i in `find tests -name '*.pb.psql' -not -path 'tests/references/*' | sort`; do \
		rel=$${i#tests/}; \
		diff tests/references/$$rel $$i; \
	done

.PHONY: test-buf-generate
test-buf-generate: build
	@mkdir -p tests/buf/proto/psql && cp -f $(NAME)/$(NAME).proto tests/buf/proto/psql/$(NAME).proto
	@sed 's#PSQL_BIN_PLACEHOLDER#$(shell pwd)/bin/protoc-gen-$(NAME)#' tests/buf/buf.gen.yaml > tests/buf/buf.gen.local.yaml
	@cd tests/buf && rm -rf gen && buf generate proto --template buf.gen.local.yaml 2>buf.stderr.txt; status=$$?; \
		cat buf.stderr.txt; \
		if grep -q "Unable to find an extension tableType" buf.stderr.txt; then \
			echo "FAIL: skip-message noise leaked to stderr (Issue 1 regression)"; exit 1; \
		fi; \
		test $$status -eq 0 || { echo "FAIL: buf generate exited $$status"; exit $$status; }
	# Checking diff between generated file and reference file (empty == identical)
	@for i in `find tests/buf/gen -name '*.pb.psql' | sort`; do \
		rel=$${i#tests/buf/}; \
		diff tests/buf/references/$$rel $$i || { echo "FAIL: $$i differs from reference (source-relative regression)"; exit 1; }; \
	done
	@test -f tests/buf/gen/demo/v1/10_tables_widget.pb.psql || { echo "FAIL: expected nested source-relative output missing"; exit 1; }
	@rm -f tests/buf/buf.gen.local.yaml tests/buf/buf.stderr.txt
	@echo "test-buf-generate OK"


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
test: build test-generate test-buf-generate test-integration


.PHONY: clean
clean:
	@rm -fv tests/*.psql
	@rm -rf tests/buf/gen tests/buf/buf.gen.local.yaml tests/buf/buf.stderr.txt tests/buf/proto/psql


.PHONY: distclean
distclean: clean
	@rm -fv bin/protoc-gen-go bin/protoc-gen-$(NAME) $(NAME)/$(NAME).pb.go
