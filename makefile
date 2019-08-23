default: test

test: install
	cd ./app && make test
	cd ./chaincode && go test .
	cd ./chaincodeSecurity && go test .

demo:
	cd ./app && make demo

install:
	go get -d -v ./...
	cd ./app && make install
	cd ./chaincode && go build
	cd ./chaincodeSecurity && go build
	go build

.PHONY: default install test demo