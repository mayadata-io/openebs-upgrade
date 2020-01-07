# Build the operator binary
FROM golang:1.13 as builder

WORKDIR /workspace

# copy build manifests
COPY Makefile Makefile

# copy go modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# ensure vendoring is up-to-date
# cache deps before building and copying source so that we don't need to
# re-download as much and so that source changes don't invalidate our
# downloaded layer
RUN make vendor

# copy go source code
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/
COPY k8s/ ./k8s/
COPY types/ ./types/
COPY controller/ ./controller/

# build the binary
RUN make bin

# Use debian as minimal base image to package the final binary
FROM debian:stretch-slim

WORKDIR /
RUN apt-get update && \
  apt-get install --no-install-recommends -y ca-certificates && \
  rm -rf /var/lib/apt/lists/*
COPY config/metac.yaml /etc/config/metac/metac.yaml
COPY templates/ /templates
COPY --from=builder /workspace/openebs-operator .

ENTRYPOINT ["/openebs-operator"]
