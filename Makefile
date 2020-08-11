.PHONY: clean clean-deps build deps run test-build docs test

clean-deps:
	go mod tidy

deps:
	go mod download

test-build: build-single
	sudo cp build/ops /usr/local/bin/ops

build:
	CGO_ENABLED=0 go build -o build/ops -ldflags "-X 'github.com/181192/ops-cli/pkg/util/version.Version=$$(git describe --tags --abbrev=0)' -X 'github.com/181192/ops-cli/pkg/util/version.GitCommit=$$(git rev-parse --short HEAD)'"

build-all:
	for arch in amd64; do \
			for os in linux darwin windows; do \
				ext=""; \
				if [ $$os = "windows" ]; then ext=".exe"; fi; \
				CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -o "build/ops_cli_"$$os"_$$arch$$ext" $(LDFLAGS) -ldflags "-X 'github.com/181192/ops-cli/pkg/util/version.Version=$$(git describe --tags --abbrev=0)' -X 'github.com/181192/ops-cli/pkg/util/version.GitCommit=$$(git rev-parse --short HEAD)'"; \
			done; \
		done;
		for arch in arm arm64; do \
			for os in linux; do \
				CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -o "build/ops_cli_"$$os"_$$arch" $(LDFLAGS) -ldflags "-X 'github.com/181192/ops-cli/pkg/util/version.Version=$$(git describe --tags --abbrev=0)' -X 'github.com/181192/ops-cli/pkg/util/version.GitCommit=$$(git rev-parse --short HEAD)'"; \
			done; \
	done;

clean: clean-deps
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
	rm -rf docs/*.md
	go run docs/docs.go

generate:
	go run pkg/codegen/main.go

generate-cleanup:
	go run pkg/codegen/cleanup/main.go
