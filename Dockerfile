#
# Copyright 2024 tofuutils authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
FROM golang:1.21 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GIT_TERMINAL_PROMPT=1

COPY ./cmd ${GOPATH}/src/github.com/tofuutils/tenv/cmd
COPY ./config ${GOPATH}/src/github.com/tofuutils/tenv/config
COPY ./pkg ${GOPATH}/src/github.com/tofuutils/tenv/pkg
COPY ./versionmanager ${GOPATH}/src/github.com/tofuutils/tenv/versionmanager
COPY ./go.mod ./go.sum ${GOPATH}/src/github.com/tofuutils/tenv/

WORKDIR ${GOPATH}/src/github.com/tofuutils/tenv
RUN go get -u ./cmd/atmos \
    && go get -u ./cmd/tenv \
    && go get -u ./cmd/terraform \
    && go get -u ./cmd/terragrunt \
    && go get -u ./cmd/tf \
    && go get -u ./cmd/tofu \
    && go mod tidy

RUN go build -ldflags="-s -w" -o atmos ./cmd/atmos \
    && go build -ldflags="-s -w" -o tenv ./cmd/tenv \
    && go build -ldflags="-s -w" -o terraform ./cmd/terraform \
    && go build -ldflags="-s -w" -o terragrunt ./cmd/terragrunt \
    && go build -ldflags="-s -w" -o tf ./cmd/tf \
    && go build -ldflags="-s -w" -o tofu ./cmd/tofu

FROM alpine:3.20
LABEL maintainer="TofuUtils Core Team"

RUN apk add --no-cache git bash openssh

COPY --from=builder go/src/github.com/tofuutils/tenv/atmos /usr/local/bin/atmos
COPY --from=builder go/src/github.com/tofuutils/tenv/tenv /usr/local/bin/tenv
COPY --from=builder go/src/github.com/tofuutils/tenv/terraform /usr/local/bin/terraform
COPY --from=builder go/src/github.com/tofuutils/tenv/terragrunt /usr/local/bin/terragrunt
COPY --from=builder go/src/github.com/tofuutils/tenv/tf /usr/local/bin/tf
COPY --from=builder go/src/github.com/tofuutils/tenv/tofu /usr/local/bin/tofu

ENTRYPOINT ["/usr/local/bin/tenv"]
