PACKAGE_NAME = github.com/bamboo-firewall/be

VERSION ?= $(shell git describe --abbrev=0 --tags)
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
BUILD_TIME ?= $(shell date +%Y-%m-%dT%H:%M:%S%z)

ORGANIZATION = ATAOCloud

LDFLAGS = -s -w -X $(PACKAGE_NAME)/buildinfo.Version=$(VERSION) \
		 -X $(PACKAGE_NAME)/buildinfo.GitBranch=$(BRANCH) \
		 -X $(PACKAGE_NAME)/buildinfo.BuildDate=$(BUILD_TIME) \
		 -X $(PACKAGE_NAME)/buildinfo.Organization=$(ORGANIZATION)

GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0

BUILD := go build -buildvcs=false -a -installsuffix cgo

.PHONY: build
build-bbfwcli:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) $(BUILD) -ldflags="$(LDFLAGS)" -o bbfwcli ./cmd/bamboofwcli

.PHONY: clean
clean:
	rm -f bbfwcli