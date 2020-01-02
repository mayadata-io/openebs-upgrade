# Copyright 2019 The MayaData Authors
# Copyright 2018 Uber Technologies, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PWD := ${CURDIR}

PACKAGE_NAME = github.com/mayadata-io/openebs-operator
PACKAGE_VERSION ?= $(shell git describe --always --tags)
OS = $(shell uname)

ALL_SRC = $(shell find . -name "*.go" | grep -v -e vendor \
	-e ".*/\..*" \
	-e ".*/_.*" \
	-e ".*/mocks.*" \
	-e ".*/*.pb.go")
ALL_PKGS = $(shell go list $(sort $(dir $(ALL_SRC))) | grep -v vendor)
ALL_PKG_PATHS = $(shell go list -f '{{.Dir}}' ./...)
FMT_SRC = $(shell echo "$(ALL_SRC)" | tr ' ' '\n')

# External tools required while building this binary or 
# to test source code, artifacts in this project
EXT_TOOLS =\
	github.com/golangci/golangci-lint/cmd/golangci-lint \
	github.com/axw/gocov/gocov \
	github.com/AlekSi/gocov-xml \
	github.com/matm/gocov-html
EXT_TOOLS_DIR = ext-tools/$(OS)

BUILD_LDFLAGS = -X $(PACKAGE_NAME)/lib/utils/build.Hash=$(PACKAGE_VERSION)
GO_FLAGS = -gcflags '-N -l' -ldflags "$(BUILD_LDFLAGS)"
GO_VERSION = 1.12

REGISTRY ?= quay.io/openebs
IMG_NAME ?= openebs-operator

### Targets to compile openebs-operator binaries
.PHONY: bins lbins
bins: $(IMG_NAME)

$(IMG_NAME): $(ALL_SRC)
	go build -tags bins $(GO_FLAGS) -o $@ cmd/main.go

### linux based binary
lbins: $(IMG_NAME).linux

$(IMG_NAME).linux: $(ALL_SRC)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on \
		go build -tags bins $(GO_FLAGS) -o $@ cmd/main.go

$(ALL_SRC): ;

### download modules to local cache
vendor: go.mod go.sum
	@go mod download

ext-tools: $(EXT_TOOLS)

### install go based tools
.PHONY: $(EXT_TOOLS)
$(EXT_TOOLS):
	@echo "Installing external tool $@"
	@GO111MODULES=on go get -u $@


### Target to build the openebs operator docker images
.PHONY: images publish
images:
	docker build -t $(REGISTRY)/$(IMG_NAME):$(PACKAGE_VERSION) -f Dockerfile .
	docker tag $(REGISTRY)/$(IMG_NAME):$(PACKAGE_VERSION) $(IMG_NAME):$(PACKAGE_VERSION)
	docker build -t $(REGISTRY)/$(IMG_NAME)-alpine:$(PACKAGE_VERSION) -f Dockerfile.alpine .
	docker tag $(REGISTRY)/$(IMG_NAME)-alpine:$(PACKAGE_VERSION) $(IMG_NAME)-alpine:$(PACKAGE_VERSION)

publish: images
	docker push $(REGISTRY)/$(IMG_NAME):$(PACKAGE_VERSION)
	docker push $(REGISTRY)/$(IMG_NAME)-alpine:$(PACKAGE_VERSION)


### Targets to test the codebase.
.PHONY: test unit-test
test: unit-test

unit-test: $(ALL_SRC) vendor ext-tools
	$(EXT_TOOLS_DIR)/gocov test $(ALL_PKGS) --tags "unit" | $(EXT_TOOLS_DIR)/gocov report

gofmt:
	@go fmt $(ALL_PKG_PATHS)

lint: ext-tools gofmt
	@echo "Running golangci-lint"
	@golangci-lint run --disable-all \
		--deadline 5m \
		--enable=misspell \
		--enable=structcheck \
		--enable=golint \
		--enable=deadcode \
		--enable=errcheck \
		--enable=varcheck \
		--enable=goconst \
		--enable=unparam \
		--enable=ineffassign \
		--enable=nakedret \
		--enable=interfacer \
		--enable=misspell \
		--enable=gocyclo \
		--enable=lll \
		--enable=dupl \
		--enable=goimports \
		$(ALL_PKG_PATHS)
