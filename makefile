default: test

test: install
	go get -d -v ./...
	cd ./app && make test
	cd ./chaincode && go build && go test .
	cd ./chaincodeSecurity && go build && go test .

demo: install
	cd ./app && make demo

install:
	#loading bug fixes to sdk packages
	cd ${GOPATH}/src/github.com/hyperledger/ \
	&& rm -rf fabric-sdk-go \
	&& git clone https://github.com/Miosolo/fabric-sdk-go

	#building app
	go build

.PHONY: default install test demo