SRC=$(shell find . -type f -name '*.go')
V=$(shell git describe --always --tags --dirty)
APP=bar
REPO=github.com/akaspin/${APP}

.PHONY: clean

all: \
	dist/${APP}-${V}-darwin-amd64.tar.gz \
	dist/${APP}-${V}-linux-amd64.tar.gz

clean:
	@rm -rf dist

dist/%/${APP}: ${SRC}
	@mkdir -p ${@D}
	CGO_ENABLED=0 GOOS=$* go build \
		-a -installsuffix cgo \
		-ldflags '-s -X main.Version=${V}' -o $@ ${REPO}

dist/${APP}-${V}-%-amd64.tar.gz: dist/%/${APP}
	tar -czf $@ -C ${<D} .

install:
	@CGO_ENABLED=0 go install \
		-a -installsuffix cgo \
		-ldflags '-s -X main.Version=${V}' ${REPO}

