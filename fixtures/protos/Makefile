
.PHONY: gen-thrift

t10/srv/constants.go:
	mkdir -p ${@D}
	thrift -out t10 --gen \
    		go:package_prefix=github.com/akaspin/bar/fixtures/protos/t10,thrift_import=github.com/apache/thrift/lib/go/thrift \
    		srv.thrift
