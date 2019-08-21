default: test

install:
	#loading bug fixes to sdk packages
	cd ${GOPATH}/src/github.com/hyperledger/ \
	&& rm -rf fabric-sdk-go \
	&& git clone https://github.com/Miosolo/fabric-sdk-go

	#building app
	go build

test: install
	go test -timeout 120s .

demo: install
	go test -v -timeout 120s .

.PHONY: default install test demo