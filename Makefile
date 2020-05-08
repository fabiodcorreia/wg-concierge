GIT_COMMIT 	= $(shell git rev-parse HEAD)
GIT_SHA    	= $(shell git rev-parse --short HEAD)
GIT_TAG    	= $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  	= $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

CMD_MAIN    = cmd/wg-concierge/main.go

ifdef BIN_VERSION
	VERSION = $(BIN_VERSION)
else
	VERSION = $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo v0)
endif

format:
	go mod tidy
	go vet ./...
	gofmt -s -w cmd/**/*.go 

run: format
	go run $(CMD_MAIN)  $(filter-out $@,$(MAKECMDGOALS))

build: format build-linux-arm build-linux-arm64 build-linux-386 build-linux-amd64
	vagrant rsync

build-linux-arm:
	GOOS=linux GOARCH=arm go build -tags release -ldflags="-s -w -X 'main.version=$(VERSION)'" -o ./build/wg-concierge_linux_arm $(CMD_MAIN)

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -tags release -ldflags="-s -w -X 'main.version=$(VERSION)'" -o ./build/wg-concierge_linux_arm64 $(CMD_MAIN)

build-linux-386:
	GOOS=linux GOARCH=386 go build -tags release -ldflags="-s -w -X 'main.version=$(VERSION)'" -o ./build/wg-concierge_linux_386 $(CMD_MAIN)

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -tags release -ldflags="-s -w -X 'main.version=$(VERSION)'" -o ./build/wg-concierge_linux_amd64 $(CMD_MAIN)

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -tags release -ldflags="-s -w -X 'main.build=$(VERSION)'" -o ./build/wg-concierge_darwin_amd64 $(CMD_MAIN)

build-dev: format build-linux-amd64
	vagrant rsync

security:
	~/go/bin/gosec ./...

clean:
	go clean
	rm -f build/wg-concierge*


.PHONY: format build clean build_linux_arm build-linux-arm64 build-linux-386 build-linux-amd64 build-darwin-amd64 build-dev security