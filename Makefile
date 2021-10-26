CC=go
FMT=gofmt
NAME=go-lnmetrics
BASE_DIR=/script
OS=linux
ARCH=386
ARM=

default: fmt lint
	$(CC) build -o $(NAME) cmd/go-lnmetrics.reporter/main.go

fmt:
	$(CC) fmt ./...

lint:
	golangci-lint run

check:
	$(CC) test -v ./...

build:
	env GOOS=$(OS) GOARCH=$(ARCH) GOARM=$(ARM) $(CC) build -o $(NAME)-$(OS)-$(ARCH) cmd/go-lnmetrics.reporter/main.go

update_utils:
	$(CC) get -u github.com/OpenLNMetrics/lnmetrics.utils
