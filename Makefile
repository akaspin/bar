SRC=$(shell find . -type f \(  -iname '*.go' ! -iname "*_test.go" \))
SRC_TEST=$(shell find . -type f -name '*_test.go')
V=$(shell git describe --always --tags --dirty)
REPO=github.com/akaspin/bar
GOOPTS=-a -installsuffix cgo -ldflags '-s -X github.com/akaspin/bar/command.Version=${V}'
#GOOPTS=-a -installsuffix cgo -ldflags '-s -X bar/command/version.Version=${V}'

HOSTNAME=$(shell ifconfig | grep 'inet ' | grep -v '127.0.0.1' | head -n1 | awk '{print $$2}')

ifdef GOBIN
	INSTALL_DIR=${GOBIN}
else
    INSTALL_DIR=${GOPATH}/bin
endif

BENCH=.
TESTS=.


.PHONY: clean uninstall test bench


all: install bindir dist

dist: dist-win dist-linux dist-darwin dist/bar-bindir-${V}.tar.gz

dist-win: dist/bar-${V}-windows-amd64.zip

dist-linux: dist/bar-${V}-linux-amd64.tar.gz

dist-darwin: dist/bar-${V}-darwin-amd64.tar.gz

distclean:
	rm -rf dist

clean:
	-find . -type d -name testdata* -exec rm -rf '{}' ';'

dist/bar-${V}-windows-amd64.zip: dist/win/bar.exe
	zip -r -j -D $@ ${<D}

dist/bar-bindir-${V}.tar.gz: dist/bindir/windows dist/bindir/linux dist/bindir/darwin
	tar -czf $@ -C ${<D} .

dist/bar-${V}-%-amd64.tar.gz: dist/%/bar
	tar -czf $@ -C ${<D} .

dist/win/bar.exe: dist/windows/bar
	@mkdir -p ${@D}
	cp dist/windows/bar dist/win/bar.exe

dist/%/bar: stubs ${SRC}
	@mkdir -p ${@D}
	CGO_ENABLED=0 GOOS=$* go build ${GOOPTS} -o $@ ${REPO}

dist/bindir/%: dist/%/bar
	@mkdir -p ${@D}
	cp dist/$*/bar dist/bindir/$*

bindir: dist/bindir/windows #dist/bindir/linux dist/bindir/darwin

install: ${INSTALL_DIR}/bar

uninstall:
	-rm ${INSTALL_DIR}/bar

${INSTALL_DIR}/bar: ${SRC}
	go install -ldflags '-s -X github.com/akaspin/bar/command.Version=${V}' ${REPO}

run-server: install dist/bindir/windows
	bar server run --log-level=DEBUG \
		--bind-http=:3000 \
		--bind-rpc=:3001 \
		--storage=block:root=testdata \
		--endpoint=${HOSTNAME}:3001 \
		--endpoint-http=http://${HOSTNAME}:3000/v1 \
		--bin-dir=dist/bindir

bench-mem:
	go test -run=XXX -bench=${BENCH} -benchmem ./...

bench:
	go test -run=XXX -bench=${BENCH} ./...

test:
	CGO_ENABLED=0 go test -run=${TESTS} ./...

test-short:
	CGO_ENABLED=0 go test -run=${TESTS} -short ./...

stubs: proto/wire/ttypes.go

clean-stubs:
	rm -rf proto/wire

proto/wire/ttypes.go: proto/proto.thrift
	mkdir -p ${@D}
	thrift -strict -v -out proto --gen \
		go:package=wire,package_prefix=github.com/akaspin/bar/proto,thrift_import=github.com/apache/thrift/lib/go/thrift,ignore_initialisms \
		$<
	rm -rf ${@D}/bar-remote

docs: doc/proto.html
	-rm -f doc/index.html

doc/proto.html: proto/proto.thrift
	thrift -out ./doc --gen html:standalone proto/proto.thrift
