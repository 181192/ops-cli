.PHONY: clean build deps run test-build docs test

deps:
	go mod tidy

test-build:
	CGO_ENABLED=0 go build -o build/ops -ldflags "-X 'github.com/181192/ops-cli/cmd.version=v0.1.0' -X 'github.com/181192/ops-cli/cmd.gitCommit=$$(git rev-parse --short HEAD)'"
	sudo cp build/ops /usr/local/bin/ops

build:
	for arch in amd64; do \
			for os in linux darwin windows; do \
				CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -o "build/ops_cli_"$$os"_$$arch" $(LDFLAGS) -ldflags "-X 'github.com/181192/ops-cli/cmd.version=v0.1.0' -X 'github.com/181192/ops-cli/cmd.gitCommit=$$(git rev-parse --short HEAD)'"; \
			done; \
		done;
		for arch in arm arm64; do \
			for os in linux; do \
				CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -o "build/ops_cli_"$$os"_$$arch" $(LDFLAGS) -ldflags "-X 'github.com/181192/ops-cli/cmd.version=v0.1.0' -X 'github.com/181192/ops-cli/cmd.gitCommit=$$(git rev-parse --short HEAD)'"; \
			done; \
	done;

clean:
	go clean
	rm -rf ./build

run: build
	./build/ops_cli_linux_amd64

test:
	go test -v ./...

docker-build:
	docker build -t ops-cli .

docker-run:
	docker run --rm ops-cli:latest version

docs:
	go run docs/docs.go
