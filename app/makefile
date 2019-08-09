PROJECT="cli-app"
BINARY="cli-app"
VERSION=1.0.0

default: install

install:
	#loading bug fixes to sdk packages
	cd ${GOPATH}/src/github.com/hyperledger/ \
	&& rm -rf fabric-sdk-go \
	&& git clone https://github.com/Miosolo/fabric-sdk-go

	#building app
	go build

.PHONY: default install test