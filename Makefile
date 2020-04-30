GIT_COMMIT 	= $(shell git rev-parse HEAD)
GIT_SHA    	= $(shell git rev-parse --short HEAD)
GIT_TAG    	= $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  	= $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

ifdef BIN_VERSION
	VERSION = $(BIN_VERSION)
else
	VERSION = $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo v0)
endif

format:
	gofmt -s -w *.go
	go mod tidy

build: format build-linux-arm build-linux-arm64 build-linux-386 build-linux-amd64
	vagrant rsync

build-linux-arm:
	GOOS=linux GOARCH=arm go build -tags release -ldflags="-s -w -X 'main.version=$(VERSION)'" -o ./build/wg-concierge_linux_arm

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -tags release -ldflags="-s -w -X 'main.version=$(VERSION)'" -o ./build/wg-concierge_linux_arm64

build-linux-386:
	GOOS=linux GOARCH=386 go build -tags release -ldflags="-s -w -X 'main.version=$(VERSION)'" -o ./build/wg-concierge_linux_386

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -tags release -ldflags="-s -w -X 'main.version=$(VERSION)'" -o ./build/wg-concierge_linux_amd64

build-dev: format build-linux-amd64
	vagrant rsync


clean:
	go clean
	rm -f build/wg-concierge*


.PHONY: format build clean build_linux_arm