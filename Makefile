CC=go
FMT=gofmt
NAME=go-lnmetrics
BASE_DIR=/script
OS=linux
ARCH=386

default: fmt
	$(CC) build -o $(NAME) cmd/go-metrics-reported/main.go

fmt:
	$(CC) fmt ./...

check:
	$(CC) test ./...

build:
	env GOOS=$(OS) GOARCH=$(ARCH) $(CC) build -o $(NAME)-$(OS)-$(ARCH) cmd/go-metrics-reported/main.go
