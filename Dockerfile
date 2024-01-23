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
RUN go get ./
RUN go build -ldflags="-s -w" -o tenv .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder go/src/github.com/tofuutils/tenv/tenv /app/
WORKDIR /app
ENTRYPOINT ["/app/tenv"]