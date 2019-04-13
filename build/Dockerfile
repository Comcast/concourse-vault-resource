FROM golang:1.11.5 AS builder

# setup build arguments
ARG OUTPUT_DIR
ARG VAULT_ADDR
ARG VAULT_TOKEN

# make the output directory
RUN mkdir -p {$OUTPUT_DIR,/go/src/github.com/comcast/concourse-vault-resource}

# copy the local code to the container location
COPY . /go/src/github.com/comcast/concourse-vault-resource

# build the binaries
WORKDIR /go/src/github.com/comcast/concourse-vault-resource

# set environment variables
ENV GO111MODULE=off CGO_ENABLED=1 GOOS=linux

# build the binaries
RUN make build 

# grab the binaries from build
FROM busybox AS resource
ARG OUTPUT_DIR
COPY --from=builder $OUTPUT_DIR/ $OUTPUT_DIR/
RUN chmod +x $OUTPUT_DIR/* && ls -lh $OUTPUT_DIR

# final resource container
FROM resource
