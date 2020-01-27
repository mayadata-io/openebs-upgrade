# Build the openebs-upgrade binary
FROM golang:1.13.5 as builder

WORKDIR /mayadata.io/openebs-upgrade/

# copy build manifests
COPY Makefile Makefile

# copy go modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# ensure vendoring is up-to-date
# cache deps before building and copying source so that we don't need to
# re-download as much and so that source changes don't invalidate our
# downloaded layer
RUN go mod download

# copy go source code
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY k8s/ k8s/
COPY types/ types/
COPY controller/ controller/

# build the binary
RUN make openebs-upgrade

# ---------------------------
# Use distroless as minimal base image to package the final binary
# Refer https://github.com/GoogleContainerTools/distroless
# ---------------------------
FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY config/metac.yaml /etc/config/metac/metac.yaml
COPY templates/ /templates
COPY --from=builder /mayadata.io/openebs-upgrade/openebs-upgrade /usr/bin

USER nonroot:nonroot

ENTRYPOINT ["/usr/bin/openebs-upgrade"]
