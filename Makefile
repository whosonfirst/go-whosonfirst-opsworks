CWD=$(shell pwd)
GOPATH := $(CWD)

build:	rmdeps deps fmt bin

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src; then rm -rf src; fi
	# mkdir -p src/github.com/whosonfirst/go-whosonfirst-opsworks/
	if test ! -d src; then mkdir src; fi
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:   rmdeps
	@GOPATH=$(GOPATH) go get -u "github.com/aws/aws-sdk-go"

vendor-deps: deps
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go

bin:	self
	@GOPATH=$(shell pwd) go build -o bin/opswof-list-instances cmd/opswof-list-instances.go
