// 1. package 정의
package main

// 2. 외부모듈 포함
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)
// 3. 체인코드 클래스 정의
type SimpleAsset struct {

}
// 4. 구조체 정의
type Asset struct {
	Key string `json:"key"`
	Value string `json:"value"`
}
type Response struct {
	Result string `json:"result"`
	Message string `json:"message"`
}
// 5. Init 함수 정의
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Init function")
	args := stub.GetStringArgs()
	if len(args) != 2{
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}
	fmt.Println("Init PutState: "+args[0]+"-"+args[1])

	var data = Asset{Key:args[0], Value:args[1]}
	dataAsBytes, _ := json.Marshal(data)
	err := stub.PutState(args[0], dataAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return shim.Success([]byte("Init process was done."))
}
// 6. Invoke 함수 정의
func (t * SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fn == "set" {
		result, err = set(stub, args)
	} else if fn == "get" {
		result, err = get(stub, args)
	} else if fn == "getHistoryForKey" {
		result, err = getHistoryForKey(stub, args)
	} else if fn== "getAllKeys" {
		result, err = getAllKeys(stub)
	} else {
		return shim.Error(err.Error())
	}
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(result))
}
// 7. set 함수 정의
func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	fmt.Println("set PutState: "+args[0]+"-"+args[1])

	var data = Asset{Key:args[0], Value:args[1]}
	dataAsBytes, _ := json.Marshal(data)
	err := stub.PutState(args[0], dataAsBytes)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return args[1], nil
}

// 8. get 함수 정의
func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	fmt.Println("get GetState: "+args[0])

	value, err := stub.GetState(args[0])
	if err != nil { // GetState 함수가 오류 난 경우
		return "", fmt.Errorf(fmt.Sprintf("Failed to create asset: %s", args[0]), err)
	}
	if value == nil { // key에 대한 값을 찾을 수 없는 경우
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	return string(value), nil
}

func getHistoryForKey(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	// 전달 인자 체크 args ? key
	if len(args) < 1 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting 1")
	}

	assetName := args[0]

	// getHistoryForKey 함수가 시작되는 로그를 남기기
	fmt.Printf("- start getHistoryForMarble: %s\n", assetName)

	resultsIterator, err := stub.GetHistoryForKey(assetName)
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for key
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		// resultsIterator에서 next 아이템 꺼내서 response 할당
		response, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForKey returning:\n%s\n", buffer.String())

	return (string)(buffer.Bytes()), nil
}

func getAllKeys(APIstub shim.ChaincodeStubInterface) (string, error) {

	startKey := "a"
	endKey := "z"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return "",err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "",err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllCars:\n%s\n", buffer.String())

	return (string)(buffer.Bytes()),nil
}

// 9. main 함수 정의
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}