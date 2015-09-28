SRC=$(shell find . -type f -name '*.go')
V=$(shell git describe --always --tags --dirty)
APP=bar
REPO=github.com/akaspin/${APP}
INSTALL_DIR=${GOPATH}/bin

.PHONY: clean uninstall

clean:
	@rm -rf dist
	@rm -rf testdata

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

