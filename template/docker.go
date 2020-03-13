package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplDocker = `FROM golang:1.13.3

# Set Arguments
ARG APP_NAME=platform

# Set go bin which doesn't appear to be set already.
ENV GOBIN /go/bin

# build directories
ENV SRC_DIR="/go/src/{{ .Src }}/{{ ucdown (getFirstService .Services).Name }}"
RUN mkdir /app
RUN mkdir -p $SRC_DIR

# Copy current directory
COPY ./Gopkg.lock ./Gopkg.toml $SRC_DIR/
WORKDIR $SRC_DIR

# Go dep!
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure -vendor-only

# Build my app
COPY . $SRC_DIR/
RUN go build -o /app/main .
CMD ["/app/main"]
`

var ListDocker = lib.List{
	FileType:     "Dockerfile",
	Template:     tmplDocker,
	Location:     "./",
	Lang:         "docker",
	ReplaceQuote: false,
}
