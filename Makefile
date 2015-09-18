SRC=$(shell find . -type f -name '*.go')
V=$(shell git describe --always --tags --dirty)
APP=bar
REPO=github.com/akaspin/${APP}

.PHONY: clean

clean:
	@rm -rf dist

install:
	@CGO_ENABLED=0 go install \
		-a -installsuffix cgo \
		-ldflags '-s -X main.Version=${V}' ${REPO}/barctl

