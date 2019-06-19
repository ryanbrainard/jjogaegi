# ########################################################## #
# Makefile for Golang Project
# Includes cross-compiling, installation, cleanup
# ########################################################## #

.PHONY: check clean install generate build build_cli build_gprc dist all

# Check for required command tools to build or stop immediately
EXECUTABLES = git go find pwd
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

BINARY=jjogaegi
PLATFORMS=darwin linux windows
VERSION=$(shell git describe --match 'v[0-9]*' --tags)
ARCHITECTURES=386 amd64
OUTPUT_DIR=dist

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

default: build

all: clean test build install

generate:
	protoc -I grpc/proto/ grpc/proto/services.proto --go_out=plugins=grpc:grpc/go/jjogaegigprc
	grpc_tools_ruby_protoc -I grpc/proto/ --ruby_out=grpc/ruby/lib --grpc_out=grpc/ruby/lib grpc/proto/services.proto

build: generate build_cli build_grpc

build_cli:
	go build ${LDFLAGS} -o $(OUTPUT_DIR)/$(BINARY) ./cmd/$(BINARY)

build_grpc: generate
	go build ${LDFLAGS} -o $(OUTPUT_DIR)/$(BINARY)-grpc-server ./cmd/$(BINARY)-grpc-server
	go build ${LDFLAGS} -o $(OUTPUT_DIR)/$(BINARY)-grpc-client-go ./cmd/$(BINARY)-grpc-client-go

dist: clean
	$(foreach GOOS, $(PLATFORMS),\
		$(foreach GOARCH, $(ARCHITECTURES),\
			$(shell env GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY)-$(GOOS)-$(GOARCH)/$(BINARY) ./cmd/$(BINARY) && zip --quiet --junk-paths --recurse-paths $(OUTPUT_DIR)/$(BINARY)-$(GOOS)-$(GOARCH).zip $(OUTPUT_DIR)/$(BINARY)-$(GOOS)-$(GOARCH))))

install: build
	cp $(OUTPUT_DIR)/$(BINARY) $(GOPATH)/bin

test:
	go test ./...

clean:
	rm -rf $(OUTPUT_DIR)
