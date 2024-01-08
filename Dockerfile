#
# Copyright © 2024 Alexander Sharov <kvendingoldo@gmail.com>, Nikolai Mishin <sanduku.default@gmail.com>, Anastasiia Kozlova <anastasiia.kozlova245@gmail.com>
#
FROM golang:1.21 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GIT_TERMINAL_PROMPT=1

COPY ./pkg ${GOPATH}/src/github.com/opentofuutils/tenv/pkg
COPY ./cmd ${GOPATH}/src/github.com/opentofuutils/tenv/cmd
COPY ./go.mod ./go.sum ./main.go ${GOPATH}/src/github.com/opentofuutils/tenv/
WORKDIR ${GOPATH}/src/github.com/opentofuutils/tenv
RUN go get ./
RUN go build -ldflags="-s -w" -o tenv .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder go/src/github.com/opentofuutils/tenv/tenv /app/
WORKDIR /app
ENTRYPOINT ["/app/tenv"]