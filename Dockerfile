FROM golang:1.13.5-alpine as build
ENV CGO_ENABLED=0

ARG GOOS=linux
ARG GOARCH=amd64
ARG LDFLAGS

RUN apk add --no-cache git ca-certificates openssh bash

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY pkg/flux/go.mod pkg/flux/go.mod
COPY pkg/flux/go.sum pkg/flux/go.sum

RUN go mod download

COPY . .
RUN GOOS=${GOOS} GOARCH=${GOARCH} go build \
  ${LDFLAGS} \
  -ldflags " \
  -X 'github.com/181192/ops-cli/pkg/util/version.Version=$(git describe --tags --abbrev=0)' \
  -X 'github.com/181192/ops-cli/pkg/util/version.GitCommit=$(git rev-parse --short HEAD)'" \
  -o "ops-cli" .

FROM build as test

ENV CI true

RUN go test ./...


FROM alpine:3.10

RUN addgroup -S app \
  && adduser -S -g app app \
  && apk add --no-cache ca-certificates

WORKDIR /home/app

COPY --from=build /app/ops-cli ops-cli

RUN chown -R app:app ./
USER app

ENTRYPOINT ["./ops-cli"]
