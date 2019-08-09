PROJECT="cli-app"
BINARY="cli-app"
VERSION=1.0.0

default: install

install:
	#copying bug fixes to sdk packages
	cp -r workarounds/* ${GOPATH}/src/
	#building app
	go build

.PHONY: default install test