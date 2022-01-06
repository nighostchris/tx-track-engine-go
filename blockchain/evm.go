package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type EvmTransaction struct {
	From        string `json:"from"`
	GasPrice    string `json:"gasPrice"`
	FeeCurrency string `json:"feeCurrency"`
	Hash        string `json:"hash"`
	Input       string `json:"input"`
	Nonce       string `json:"nonce"`
	To          string `json:"to"`
	Value       string `json:"value"`
}

type EvmBlock struct {
	Hash         string           `json:"hash"`
	Number       string           `json:"number"`
	Timestamp    string           `json:"timestamp"`
	Transactions []EvmTransaction `json:"transactions"`
}

type EvmRpcRequest struct {
	Id      string        `json:"id"`
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type EvmGetLatestBlockResponse struct {
	Id      string `json:"id"`
	JsonRpc string `json:"jsonrpc"`
	Result  string `json:"result"`
}

type EvmGetBlockByNumberResponse struct {
	Id      string   `json:"id"`
	JsonRpc string   `json:"jsonrpc"`
	Result  EvmBlock `json:"result"`
}

func hexToDec(input string) (decStr string, err error) {
	postProcessed := strings.Replace(input, "0x", "", -1)
	decimal, decParseError := strconv.ParseUint(postProcessed, 16, 64)

	if decParseError != nil {
		return "", decParseError
	}

	return strconv.FormatUint(decimal, 10), nil
}

func GetLatestBlock(nodeUrl string) (block string, err error) {
	const logHead = "[EVM Node RPC] GetLatestBlock -"

	params, _ := json.Marshal(&EvmRpcRequest{
		JsonRpc: "2.0",
		Method:  "eth_blockNumber",
		Params:  []interface{}{},
		Id:      "1",
	})

	response, rpcError := http.Post(nodeUrl, "application/json", bytes.NewBuffer(params))

	if rpcError != nil {
		return "", fmt.Errorf("%s %s", logHead, rpcError.Error())
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("%s %s", logHead, response.Status)
	} else {
		var data EvmGetLatestBlockResponse

		decodeErr := json.NewDecoder(response.Body).Decode(&data)

		if decodeErr != nil {
			return "", fmt.Errorf("%s %s", logHead, decodeErr.Error())
		}

		decBlock, _ := hexToDec(data.Result)

		fmt.Printf("%s %s (%s)\n", logHead, decBlock, data.Result)
		return data.Result, nil
	}
}

func GetBlockByNumber(nodeUrl string, number string) (block EvmBlock, err error) {
	const logHead = "[EVM Node RPC] GetBlockByNumber - "
	decNumber, _ := hexToDec(number)

	fmt.Printf("%s%s\n", logHead, decNumber)

	params, _ := json.Marshal(&EvmRpcRequest{
		JsonRpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{number, true},
		Id:      "1",
	})

	response, rpcError := http.Post(nodeUrl, "application/json", bytes.NewBuffer(params))

	if rpcError != nil {
		return EvmBlock{}, fmt.Errorf("%s %s", logHead, rpcError.Error())
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return EvmBlock{}, fmt.Errorf("%s %s", logHead, response.Status)
	} else {
		var data EvmGetBlockByNumberResponse

		decodeErr := json.NewDecoder(response.Body).Decode(&data)

		if decodeErr != nil {
			return EvmBlock{}, fmt.Errorf("%s %s", logHead, decodeErr.Error())
		}

		fmt.Printf("%s Block #%s have %d transactions\n", logHead, decNumber, len(data.Result.Transactions))

		return data.Result, nil
	}
}

func ProcessBlock(nodeUrl string, number string, addresses []string, processed chan string) {
	const logHead = "[EVM Process Block]"
	decNumber, _ := hexToDec(number)

	block, getBlockByNumberError := GetBlockByNumber(nodeUrl, number)

	if getBlockByNumberError != nil {
		fmt.Println(getBlockByNumberError.Error())
	}

	interestedTransaction := 0

	// Check if there is any transactions containing interested address as sender or receipient
	for _, transaction := range block.Transactions {
		if Contains(addresses, transaction.From, false) || Contains(addresses, transaction.To, false) {
			fmt.Printf("%s Found interested transaction - %s", logHead, transaction.Hash)
			interestedTransaction++
		}
	}

	fmt.Printf("%s Finished processing block %s and found %d interested transaction(s)\n", logHead, decNumber, interestedTransaction)

	processed <- number
}
