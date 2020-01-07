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

REGISTRY ?= quay.io/openebs
IMG_NAME ?= openebs-operator

all: bin

### Targets to compile openebs-operator binary
.PHONY: bin
bin: vendor $(IMG_NAME)

$(IMG_NAME): fmt vet
	@echo "+ Generating $(IMG_NAME) binary"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on \
		go build -o $@ ./cmd/main.go

# go mod download modules to local cache
# make vendored copy of dependencies
# install other go binaries for code generation
.PHONY: vendor
vendor: go.mod go.sum
	@GO111MODULE=on go mod download
	@GO111MODULE=on go mod vendor

# Run tests
.PHONY: test
test: fmt vet
	@go test ./... -coverprofile cover.out

# Run go fmt against code
.PHONY: fmt
fmt:
	@go fmt ./...

# Run go vet against code
.PHONY: vet
vet:
	@go vet ./...

# Build the docker image
.PHONY: docker-build
docker-build: test
	docker build -t $(REGISTRY)/$(IMG_NAME):$(PACKAGE_VERSION) .

# Push the docker image
.PHONY: docker-push
docker-push: docker-build
	docker push $(REGISTRY)/$(IMG_NAME):$(PACKAGE_VERSION)