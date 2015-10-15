SRC=$(shell find . -type f \(  -iname '*.go' ! -iname "*_test.go" \))
SRC_TEST=$(shell find . -type f -name '*_test.go')
V=$(shell git describe --always --tags --dirty)
REPO=github.com/akaspin/bar
GOOPTS=-a -installsuffix cgo -ldflags '-s -X main.Version=${V}'

HOSTNAME=$(shell ifconfig | grep 'inet ' | grep -v '127.0.0.1' | head -n1 | awk '{print $$2}')

ifdef GOBIN
	INSTALL_DIR=${GOBIN}
else
    INSTALL_DIR=${GOPATH}/bin
endif

BENCH=.
TESTS=.


.PHONY: clean uninstall test bench


dist: dist-win dist-linux dist-darwin

dist-win: dist/bar-${V}-windows-amd64.zip

dist-linux: dist/bar-${V}-linux-amd64.tar.gz

dist-darwin: dist/bar-${V}-darwin-amd64.tar.gz

distclean:
	rm -rf dist

clean:
	-find . -type d -name testdata* -exec rm -rf '{}' ';'

dist/bar-${V}-windows-amd64.zip: dist/windows/barc.exe
	zip -r -j -D $@ ${<D}

dist/bar-${V}-windows-amd64.tar.gz: dist/windows/barc.exe
	tar -czf $@ -C ${<D} .

dist/bar-${V}-%-amd64.tar.gz: dist/%/barc dist/%/bard
	tar -czf $@ -C ${<D} .

dist/windows/%.exe: ${SRC}
	@mkdir -p ${@D}
	CGO_ENABLED=0 GOOS=windows go build ${GOOPTS} -o $@ ${REPO}/$*

dist/%/bard: ${SRC}
	@mkdir -p ${@D}
	CGO_ENABLED=0 GOOS=$* go build ${GOOPTS} -o $@ ${REPO}/$(@F)

dist/%/barc: ${SRC}
	@mkdir -p ${@D}
	CGO_ENABLED=0 GOOS=$* go build ${GOOPTS} -o $@ ${REPO}/$(@F)

install: ${INSTALL_DIR}/bard ${INSTALL_DIR}/barc

uninstall:
	rm ${INSTALL_DIR}/bard ${INSTALL_DIR}/barc

${INSTALL_DIR}/bard: ${SRC}
	CGO_ENABLED=0 go install ${GOOPTS} ${REPO}/bard

${INSTALL_DIR}/barc: ${SRC}
	CGO_ENABLED=0 go install ${GOOPTS} ${REPO}/barc

run-server: ${INSTALL_DIR}/bard dist/windows/barc.exe
	bard -log-level=DEBUG \
		-bind-http=:3000 \
		-bind-rpc=:3001 \
		-storage-block-root=testdata \
		-rpc=${HOSTNAME}:3001 \
		-http=http://${HOSTNAME}:3000/v1 \
		-barc-exe=dist/windows/barc.exe

bench-mem:
	go test -run=XXX -bench=${BENCH} -benchmem ./...

bench:
	go test -run=XXX -bench=${BENCH} ./...

test:
	CGO_ENABLED=0 go test -run=${TESTS} ./...

test-short:
	CGO_ENABLED=0 go test -run=${TESTS} -short ./...

stubs: proto/bar/ttypes.go

clean-stubs:
	rm -rf proto/wire

proto/bar/ttypes.go: proto/proto.thrift
	mkdir -p ${@D}
	thrift -strict -v -out proto --gen \
		go:package=wire,package_prefix=github.com/akaspin/bar/proto,thrift_import=github.com/apache/thrift/lib/go/thrift,ignore_initialisms \
		$<
	rm -rf ${@D}/bar-remote