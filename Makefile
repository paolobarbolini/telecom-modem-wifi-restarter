NAME = "modem-wifi-restarter"
VERSION = "v1.0"

NAME_PREFIX = "build/$(NAME)-$(VERSION)"

GO_LDFLAGS = "-s -w"
GO_FILES ?= $(shell find . -name '*.go')

all: build

install:
	go get -u gopkg.in/PuerkitoBio/goquery.v1

verify: fmt vet lint

fmt:
	go fmt . | awk '{ print "Please run go fmt"; exit 1 }'

vet:
	go tool vet *.go

lint:
	go get -u github.com/golang/lint/golint
	golint .

build: $(GO_FILES) install
	GOOS=linux GOARCH=amd64 go build -ldflags $(GO_LDFLAGS) -o $(NAME_PREFIX)_linux_amd64
	GOOS=linux GOARCH=386 go build -ldflags $(GO_LDFLAGS) -o $(NAME_PREFIX)_linux_386
	GOOS=linux GOARCH=arm GOARM=7 go build -ldflags $(GO_LDFLAGS) -o $(NAME_PREFIX)_linux_arm
	GOOS=linux GOARCH=arm64 go build -ldflags $(GO_LDFLAGS) -o $(NAME_PREFIX)_linux_arm64
	GOOS=darwin GOARCH=386 go build -ldflags $(GO_LDFLAGS) -o $(NAME_PREFIX)_darwin_386
	GOOS=darwin GOARCH=amd64 go build -ldflags $(GO_LDFLAGS) -o $(NAME_PREFIX)_darwin_amd64
	GOOS=windows GOARCH=386 go build -ldflags $(GO_LDFLAGS) -o $(NAME_PREFIX)_windows_386.exe
	GOOS=windows GOARCH=amd64 go build -ldflags $(GO_LDFLAGS) -o $(NAME_PREFIX)_windows_amd64.exe

clean:
	@rm -rf build/
