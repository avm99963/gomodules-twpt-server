FROM golang:latest@sha256:672a2286da3ee7a854c3e0a56e0838918d0dbb1c18652992930293312de898a6 as builder

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

FROM alpine@sha256:51b67269f354137895d43f3b3d810bfacd3945438e94dc5ac55fdac340352f48
LABEL maintainer "me@avm99963.com"

RUN apk add --no-cache tini
COPY --from=builder /go/bin/server /
ENTRYPOINT ["/sbin/tini", "--", "/server"]
