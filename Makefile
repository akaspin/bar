SRC=$(shell find . -type f -name '*.go')
V=$(shell git describe --always --tags --dirty)
REPO=github.com/akaspin/bar
INSTALL_DIR=${GOPATH}/bin

.PHONY: clean uninstall

all: \
	dist/bar-${V}-windows-amd64.tar.gz \
	dist/bar-${V}-linux-amd64.tar.gz \
	dist/bar-${V}-darwin-amd64.tar.gz

clean:
	@rm -rf dist
	@rm -rf testdata

dist/bar-${V}-windows-amd64.tar.gz: dist/windows/barc.exe dist/windows/bard.exe
	tar -czf $@ -C ${<D} .

dist/bar-${V}-%-amd64.tar.gz: dist/%/barc dist/%/bard
	tar -czf $@ -C ${<D} .

dist/windows/%.exe: ${SRC}
	@mkdir -p ${@D}
	CGO_ENABLED=0 GOOS=windows go build \
		-a -installsuffix cgo \
		-ldflags '-s -X main.Version=${V}' -o $@ ${REPO}/$*

dist/%/bard: ${SRC}
	@mkdir -p ${@D}
	CGO_ENABLED=0 GOOS=$* go build \
		-a -installsuffix cgo \
		-ldflags '-s -X main.Version=${V}' -o $@ ${REPO}/$(@F)

dist/%/barc: ${SRC}
	@mkdir -p ${@D}
	CGO_ENABLED=0 GOOS=$* go build \
		-a -installsuffix cgo \
		-ldflags '-s -X main.Version=${V}' -o $@ ${REPO}/$(@F)


install: ${INSTALL_DIR}/bard ${INSTALL_DIR}/barc

uninstall:
	rm ${INSTALL_DIR}/bard ${INSTALL_DIR}/barc

${INSTALL_DIR}/bard: ${SRC}
	CGO_ENABLED=0 go install \
		-a -installsuffix cgo \
		-ldflags '-s -X main.Version=${V}' ${REPO}/bard

${INSTALL_DIR}/barc: ${SRC}
	CGO_ENABLED=0 go install \
		-a -installsuffix cgo \
		-ldflags '-s -X main.Version=${V}' ${REPO}/barc

