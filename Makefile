VERSION := $(shell git describe --tags --always --dirty="-dev")
LDFLAGS := -ldflags='-X "main.version=$(VERSION)"'
MODMODE := -mod=readonly

version:
	echo ${VERSION}

test:
	go test ${MODMODE} -cover ./...

lint:
	go vet ${MODMODE} ./...

build:
	go build ${LDFLAGS} ${MODMODE} -trimpath .

.PHONY: *
