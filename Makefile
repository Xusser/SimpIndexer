export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w -X github.com/Xusser/SimpIndexer/build.BuildDate=`date -u +%Y%m%d` -X github.com/Xusser/SimpIndexer/build.BuildCommit=`git rev-parse --short HEAD`
OS := windows linux darwin openbsd
ARCH := amd64 386 arm arm64
OSARCH := !darwin/arm !darwin/386
OUTPUT := dist/simpindexer-{{.OS}}-{{.Arch}}

all: clean fmt build

build: release upx

fmt:
	go fmt ./...

release:
	gox -ldflags "$(LDFLAGS)" -os="${OS}" -arch="${ARCH}" -osarch="$(OSARCH)" -output="$(OUTPUT)" .

upx:
	upx --lzma dist/simpindexer-linux-*
	upx --lzma dist/simpindexer-windows-*
	upx --lzma dist/simpindexer-darwin-*
	
clean:
	rm -f ./dist/*