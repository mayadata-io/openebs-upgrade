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
COPY cmd/ cmd/
COPY pkg/ pkg/

# build the binary
RUN make bin

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/openebs-operator .
USER nonroot:nonroot

ENTRYPOINT ["/openebs-operator"]
