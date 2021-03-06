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

	res := stub.MockInit("1", [][]byte{[]byte("init")})
	fmt.Println("Init result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("Yongmao@ANZBank"), []byte("0")})
	fmt.Println("create Yongmao result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("add"), []byte("Yongmao@ANZBank"), []byte("10")})
	fmt.Println("add Yongmao result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("reduce"), []byte("Yongmao@ANZBank"), []byte("10")})
	fmt.Println("reduce Yongmao result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("Songyue@ANZBank"), []byte("0")})
	fmt.Println("create Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Songyue@ANZBank")})
	fmt.Println("get Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("Songyue@ANZBank"), []byte("10")})
	fmt.Println("create Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("transfer"), []byte("Yongmao@ANZBank"), []byte("Songyue@ANZBank"), []byte("10")})
	fmt.Println("transfer Yongmao Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Songyue@ANZBank")})
	fmt.Println("get Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Yongmao@ANZBank")})
	fmt.Println("get Yongmao result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("add"), []byte("Yongmao@ANZBank"), []byte("100")})
	fmt.Println("add Yongmao result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("transfer"), []byte("Yongmao@ANZBank"), []byte("Songyue@ANZBank"), []byte("10")})
	fmt.Println("transfer Yongmao Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Songyue@ANZBank")})
	fmt.Println("get Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Yongmao@ANZBank")})
	fmt.Println("get Yongmao result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("query"), []byte("in"), []byte("Songyue@ANZBank")})
	fmt.Println("query In Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("query"), []byte("out"), []byte("Yongmao@ANZBank")})
	fmt.Println("query Out Yongmao result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("Yuhao@ANZBank"), []byte("0")})
	fmt.Println("create Yuhao result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("delete"), []byte("Yuhao@ANZBank")})
	fmt.Println("delete Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Yuhao@ANZBank")})
	fmt.Println("get Yuhao result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("RollBack"), []byte("Yongmao@ANZBank"), []byte("Songyue@ANZBank"), []byte("1")})
	fmt.Println("RollBack Yongmao Songyue 1 result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Songyue@ANZBank")})
	fmt.Println("get Songyue result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Yongmao@ANZBank")})
	fmt.Println("get Yongmao result: ", string(res.Payload))

	res = stub.MockInvoke("1", [][]byte{[]byte("RollBack"), []byte("Yongmao@ANZBank"), []byte("Yuhao@ANZBank"), []byte("1")})
	fmt.Println("RollBack Yongmao Yuhao 1 result: ", string(res.Payload))
}

func BenchmarkCreateGetDelete(b *testing.B) {
	cc := new(SimpleAsset)
	stub := shim.NewMockStub("test", cc)

	for i := 0; i < b.N; i++ {
		stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("Songyue"), []byte("0")})
		stub.MockInvoke("1", [][]byte{[]byte("get"), []byte("Songyue")})
		stub.MockInvoke("1", [][]byte{[]byte("delete"), []byte("Songyue")})
	}
}

func BenchmarkCreateTransferQueryRollBack(b *testing.B) {
	cc := new(SimpleAsset)
	stub := shim.NewMockStub("test", cc)

	for i := 0; i < b.N; i++ {
		stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("Songyue"), []byte("0")})
		stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("Yongmao"), []byte("100000")})
		stub.MockInvoke("1", [][]byte{[]byte("transfer"), []byte("Yongmao"), []byte("Songyue"), []byte("1")})
		stub.MockInvoke("1", [][]byte{[]byte("query"), []byte("out"), []byte("Yongmao")})
		stub.MockInvoke("1", [][]byte{[]byte("query"), []byte("in"), []byte("Songyue")})
		stub.MockInvoke("1", [][]byte{[]byte("RollBack"), []byte("Yongmao"), []byte("Songyue"), []byte("1")})
	}
}

func BenchmarkCreateAddReduceDelete(b *testing.B) {
	cc := new(SimpleAsset)
	stub := shim.NewMockStub("test", cc)

	for i := 0; i < b.N; i++ {
		stub.MockInvoke("1", [][]byte{[]byte("create"), []byte("Songyue"), []byte("0")})
		stub.MockInvoke("1", [][]byte{[]byte("add"), []byte("Songyue"), []byte("1")})
		stub.MockInvoke("1", [][]byte{[]byte("reduce"), []byte("Songyue"), []byte("1")})
		stub.MockInvoke("1", [][]byte{[]byte("delete"), []byte("Songyue")})
	}
}
