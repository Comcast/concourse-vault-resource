# PROJECT_NAME - the name of the project
PROJECT_NAME := concourse-vault-resource

# PROJECT_DIR  - the project directory in the $GOPATH
PROJECT_DIR  := $(GOPATH)/src/github.com/comcast/$(PROJECT_NAME)

# OUTPUT_DIR   - the output directory for the binary builds. 
# also the output directory in the Dockerfile
OUTPUT_DIR   := /opt/resource

# GIT_COMMIT   - the git commit
GIT_COMMIT   := $(shell git rev-parse --short HEAD)

# VERSION      - the version as read from branch version file version
VERSION      := $(shell git show version:version)

# DOCKER_IMAGE - the docker image repository and name
DOCKER_IMAGE ?= hub.example.com/foo/concourse-vault-resource

# GOLANG SPECIFICS - the variables with ?= are sane defaults should they not
# already be set
GOVERSION   := 1.11.5
GOMAXPROCS  ?= 4
GO111MODULE ?= on
GOPATH      ?= $(shell go env GOPATH)
GOTAGS      ?=
LD_FLAGS    ?= \
	-s \
	-w \
	-extldflags "-static" \
	-X $(PROJECT_DIR)/version.Name=$(PROJECT_NAME) \
	-X $(PROJECT_DIR)/version.GitCommit=$(GIT_COMMIT)

# Concourse Vault Resource Variables
VAULT_ADDR  ?= https://vault.example.com:8200
VAULT_TOKEN := $(shell echo $$VAULT_TOKEN)

all: usage
usage: Makefile
	@echo
	@echo "$(PROJECT_NAME) supports the following:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
.PHONY: usage


## build - build the binaries
build: lint test
	@echo "Building binaries with LD_FLAGS \"$(LD_FLAGS)\""
	@for b in "check" "in" "out"; do \
		go build -a \
		  -ldflags "$(LD_FLAGS)" \
		  -o $(OUTPUT_DIR)/$$b \
		  -tags "$(GOTAGS)" \
		  cmd/$$b/$$b.go; \
	done
.PHONY: build

## deps - will download dependencies
deps:
	export GO111MODULE=$(GO111MODULE)
	go mod download
	go mod vendor

## fmt - will execute go fmt
fmt:
	go fmt ./... 

## image - will build the docker image
image: 
	docker build \
		--build-arg OUTPUT_DIR=$(OUTPUT_DIR) \
		--build-arg VAULT_ADDR=$(VAULT_ADDR)\
		--build-arg VAULT_TOKEN=$(VAULT_TOKEN) \
		-f build/Dockerfile --rm -t $(DOCKER_IMAGE) .

## install - will install the binary
install:
	go install -v

## lint - will lint the code
lint:
	@curl -s -L https://git.io/vp6lP | sh -s -- -b $(GOPATH)/bin
	gometalinter --config=.gometalinter.json \
		--fast --deadline=90s --vendor --errors ./...

## push - will push the docker image to the docker repository
push:
	docker push $(DOCKER_IMAGE) 

## test - will execute any tests
test:
	go test -v ./...
.PHONY: test 
