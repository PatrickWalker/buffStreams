############################
# STEP 1 build executable binary
############################
FROM golang:1.13-alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git && apk add --update make
WORKDIR $GOPATH/src/github.com/PatrickWalker/buffStreams
COPY . .
# Fetch dependencies.
RUN  make -f MakeFile docker
############################
# STEP 2 build a small image
############################
FROM alpine
# Copy our static executable.
COPY --from=builder /go/bin/buff-api /go/bin/buff-api
COPY config.json /go/bin/config.json
COPY migrations /go/bin/migrations
WORKDIR /go/bin/
ENTRYPOINT [ "./buff-api" ]
