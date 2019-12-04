FROM golang:1.12 as build
ENV CGO_ENABLED=0

ARG GOOS=linux
ARG GOARCH=amd64
ARG LDFLAGS

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN GOOS=${GOOS} GOARCH=${GOARCH} go build \
  ${LDFLAGS} \
  -ldflags " \
  -X 'github.com/181192/ops-cli/cmd.version=v0.1.0' \
  -X 'github.com/181192/ops-cli/cmd.gitCommit=$(git rev-parse HEAD)'" \
  -o "ops-cli" .


FROM alpine:3.10

RUN addgroup -S app \
  && adduser -S -g app app \
  && apk add --no-cache ca-certificates

WORKDIR /home/app

COPY --from=build /app/ops-cli ops-cli

RUN chown -R app:app ./
USER app

ENTRYPOINT ["./ops-cli"]
CMD ["version"]
