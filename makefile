default:
	go get -d -v ./...
	cd ./app && make install
	cd ./chaincode && go build
	cd ./chaincode-security && go build