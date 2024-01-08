FROM golang:1.18 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GIT_TERMINAL_PROMPT=1

COPY ./internal ${GOPATH}/src/github.com/kvendingoldo/aws-cognito-backup-lambda/internal
COPY ./go.mod ./go.sum ./main.go ${GOPATH}/src/github.com/kvendingoldo/aws-cognito-backup-lambda/
WORKDIR ${GOPATH}/src/github.com/kvendingoldo/aws-cognito-backup-lambda
RUN go get ./
RUN go build -ldflags="-s -w" -o lambda .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder go/src/github.com/kvendingoldo/aws-cognito-backup-lambda/lambda /app/
WORKDIR /app
ENTRYPOINT ["/app/lambda"]
