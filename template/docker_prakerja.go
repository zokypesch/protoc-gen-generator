package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplDockerPrakerja = `FROM prakerja/builder-base:go1.13.3

# Set Arguments
ARG APP_NAME={{ ucfirst (getFirstService .Services).Name }}

# Set go bin which doesn't appear to be set already.
ENV GOBIN /go/bin

# set credential
ARG SSH_PRIVATE_KEY

RUN mkdir -p ~/.ssh && umask 0077 && echo "${SSH_PRIVATE_KEY}" > ~/.ssh/id_rsa \
    && git config --global url."git@gitlab.com:".insteadOf https://gitlab.com/ \
    && ssh-keyscan gitlab.com >> ~/.ssh/known_hosts &> /dev/null

RUN export GOPRIVATE=gitlab.com/prakerja

# build directories
ENV SRC_DIR="/go/src/gitlab.com/prakerja/{{ .Src }}"
RUN mkdir /app
RUN mkdir -p $SRC_DIR

# Copy current directory
COPY ./go.mod ./go.sum $SRC_DIR/
WORKDIR $SRC_DIR

# Build my app
COPY . $SRC_DIR/

# go mod download
RUN go mod download

RUN go build -o /app/main .
CMD ["/app/main"]

`

var ListDockerPrakerja = lib.List{
	FileType:     "Dockerfile",
	Template:     tmplDockerPrakerja,
	Location:     "./",
	Lang:         "docker",
	ReplaceQuote: false,
}
