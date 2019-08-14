default:
	go get -d -v ./...
	cd ./app && make install
	cd ./chaincode && go build && go test .
	cd ./chaincodeSecurity && go build && go test .