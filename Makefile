#!/usr/bin/make -f

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

export GO111MODULE = on

# process build tags

build_tags = netgo

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/Dipper-Labs/Dipper-Protocol/version.Name=dip \
		  -X github.com/Dipper-Labs/Dipper-Protocol/version.ServerName=dipd \
		  -X github.com/Dipper-Labs/Dipper-Protocol/version.ClientName=dipcli \
		  -X github.com/Dipper-Labs/Dipper-Protocol/version.Version=$(VERSION) \
		  -X github.com/Dipper-Labs/Dipper-Protocol/version.Commit=$(COMMIT) \
		  -X "github.com/Dipper-Labs/Dipper-Protocol/version.BuildTags=$(build_tags_comma_sep)"

ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/Dipper-Labs/Dipper-Protocol/types.DBBackend=cleveldb
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

all: get_tools install

get_tools:
	cd scripts && $(MAKE) get_tools

check_dev_tools:
	cd scripts && $(MAKE) check_dev_tools

get_dev_tools:
	cd scripts && $(MAKE) get_dev_tools

### Generate swagger docs for dipd
update_dipd_swagger_docs:
	@statik -src=client/lcd/swagger-ui -dest=client/lcd -f

build: go.sum
ifeq ($(OS),Windows_NT)
	go build $(BUILD_FLAGS) -o build/dipd.exe ./cmd/dipd
	go build $(BUILD_FLAGS) -o build/dipcli.exe ./cmd/dipcli
else
	go build $(BUILD_FLAGS) -o build/dipd ./cmd/dipd
	go build $(BUILD_FLAGS) -o build/dipcli ./cmd/dipcli
endif

build-linux: go.sum
	GOOS=linux GOARCH=amd64 $(MAKE) build

install: go.sum
	go install $(BUILD_FLAGS) ./cmd/dipd
	go install $(BUILD_FLAGS) ./cmd/dipcli


########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd/dipd -d 2 | dot -Tpng -o dependency-graph.png

clean:
	rm -rf  build/



##############################################
### Test
PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation' | grep -v mock | grep -v 'Dipper-Protocol/tests' | grep -v 'Dipper-Protocol/crypto' | grep -v '/simapp')
PACKAGES_CRYPTO=$(shell go list ./... | grep '/Dipper-Labs/Dipper-Protocol/crypto')

test: test_unit_crypto test_unit
simapptest:
	@go test -mod=readonly github.com/Dipper-Labs/Dipper-Protocol/app/simapp \
        -run=TestFullAppSimulation \
        -Enabled=true \
        -NumBlocks=100 \
        -BlockSize=200 \
        -Commit=true \
        -Seed=99 \
        -Period=5 \
        -v -timeout 24h

test_unit:
	@go test -mod=readonly $(PACKAGES_NOSIMULATION)

test_unit_crypto:
	@go test -mod=readonly -tags "ledger test_ledger_mock" $(PACKAGES_CRYPTO)
