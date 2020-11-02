.PHONY: service

all: service web

service:
	go build -o bin/service cmd/service/main.go

web:
	go build -o bin/web cmd/web/main.go

fetch:
	go get -v github.com/axw/gocov/gocov
	go get -v github.com/AlekSi/gocov-xml
	go get -v golang.org/x/lint/golint
	go get -v gopkg.in/alecthomas/gometalinter.v2
	go get -v github.com/kisielk/errcheck

mock-fetch:
	go get github.com/golang/mock/gomock
	go get github.com/golang/mock/mockgen

lint:
	$(ENV) $(GOPATH)/bin/gometalinter.v2 ./pkg/...

test: fetch lint
	$(ENV) $(GOPATH)/bin/gocov test ./pkg/... | $(GOPATH)/bin/gocov-xml > coverage.xml

thrift:
	mkdir -p gen-go
	tzonecli generate-go thrift/purple.thrift -o gen-go

mock: mock-fetch
	mkdir -p pkg/basic/rpc/mock
	mockgen -source=pkg/basic/rpc/member.go -package=mock -destination=pkg/basic/rpc/mock/mock_member.go

