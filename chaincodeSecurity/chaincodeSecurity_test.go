// this file is used to have unit test of chaincode
// however, due to the difficulty of adding msp information in mockStub,
// the real unit test file can function in the version that do not have
// authority control.
package chaincodeSecurity

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestChaincode(t *testing.T) {
	cc := new(SimpleAsset)
	stub := shim.NewMockStub("test", cc)

	res := stub.MockInit("1", [][]byte{[]byte("init"), []byte("Yongmao"), []byte("100")})
	fmt.Println("init Yongmao result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("Songyue"), []byte("0")})
	fmt.Println("create Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Songyue")})
	fmt.Println("get Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("Songyue"), []byte("10")})
	fmt.Println("create Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Songyue")})
	fmt.Println("get Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("transfer"), []byte("Yongmao"), []byte("Songyue"), []byte("10")})
	fmt.Println("transfer Yongmao Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("query"), []byte("in"), []byte("Songyue")})
	fmt.Println("query In Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("query"), []byte("out"), []byte("Yongmao")})
	fmt.Println("query Out Yongmao result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("delete"), []byte("Songyue")})
	fmt.Println("delete Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Songyue")})
	fmt.Println("get Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("query"), []byte("in"), []byte("Songyue")})
	fmt.Println("query In Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("query"), []byte("out"), []byte("Songyue")})
	fmt.Println("query Out result: ", string(res.Payload))
}
