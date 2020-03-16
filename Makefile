PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/Dipper-Protocol/version.Name=dipperProtocol \
	-X github.com/Dipper-Protocol/version.ServerName=dipd \
	-X github.com/Dipper-Protocol/version.ClientName=dipcli \
	-X github.com/Dipper-Protocol/version.Version=$(VERSION) \
	-X github.com/Dipper-Protocol/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

include Makefile.ledger
all: install

install: go.sum
		go install $(BUILD_FLAGS) ./cmd/dipd
		go install $(BUILD_FLAGS) ./cmd/dipcli

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

test:
	@go test -mod=readonly $(PACKAGES)