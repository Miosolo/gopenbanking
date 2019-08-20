package chaincode

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestChaincode(t *testing.T) {
	cc := new(SimpleAsset)
	stub := shim.NewMockStub("test", cc)
	res := stub.MockInit("1", [][]byte{[]byte("init"), []byte("Yongmao"), []byte("100")})
	fmt.Println("init result: ", string(res.Payload))
	res = stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("Songyue"), []byte("0")})
	fmt.Println("create result: ", string(res.Payload))
	res = stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("Songyue"), []byte("0")})
	fmt.Println("create result: ", string(res.Payload))
	res = stub.MockInvoke("1", [][]byte{[]byte("transfer"), []byte("Yongmao"), []byte("Songyue"), []byte("10")})
	fmt.Println("transfer result: ", string(res.Payload))
	res = stub.MockInvoke("1", [][]byte{[]byte("query"), []byte("in"), []byte("Songyue")})
	fmt.Println("query result: ", string(res.Payload))
	res = stub.MockInvoke("1", [][]byte{[]byte("query"), []byte("out"), []byte("Yongmao")})
	fmt.Println("query result: ", string(res.Payload))
	res = stub.MockInvoke("1", [][]byte{[]byte("delete"), []byte("Songyue")})
	fmt.Println("delete result: ", string(res.Payload))
	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Songyue")})
	fmt.Println("get result: ", string(res.Payload))
	res = stub.MockInvoke("1", [][]byte{[]byte("query"), []byte("in"), []byte("Songyue")})
	fmt.Println("query result: ", string(res.Payload))
	res = stub.MockInvoke("1", [][]byte{[]byte("query"), []byte("out"), []byte("Songyue")})
	fmt.Println("query result: ", string(res.Payload))
}
