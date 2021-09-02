FROM golang:latest as builder

ENV CGO_ENABLED=0

ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org

RUN mkdir /gocache
ENV GOCACHE /gocache

WORKDIR /go/src/gomodules.avm99963.com/twpt-server
COPY go.mod /go/src/gomodules.avm99963.com/twpt-server/go.mod
COPY go.sum /go/src/gomodules.avm99963.com/twpt-server/go.sum
COPY cmd /go/src/gomodules.avm99963.com/twpt-server/cmd
COPY internal /go/src/gomodules.avm99963.com/twpt-server/internal
COPY api_proto /go/src/gomodules.avm99963.com/twpt-server/api_proto

WORKDIR /go/src/gomodules.avm99963.com/twpt-server/cmd/server
RUN go install gomodules.avm99963.com/twpt-server/cmd/server

FROM alpine
LABEL maintainer "me@avm99963.com"

RUN apk add --no-cache tini
COPY --from=builder /go/bin/server /
ENTRYPOINT ["/sbin/tini", "--", "/server"]
