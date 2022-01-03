package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	const nodeUrl = "https://alfajores-forno.celo-testnet.org"

	blockNumber := GetLatestBlock(nodeUrl)
	GetBlockByNumber(nodeUrl, blockNumber)
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

type EvmGetBlockByNumberResponse struct {
	Id      string   `json:"id"`
	JsonRpc string   `json:"jsonrpc"`
	Result  EvmBlock `json:"result"`
}

func GetBlockByNumber(nodeUrl string, block string) interface{} {
	const logHead = "[Celo Node RPC] GetBlockByNumber - "

	params, _ := json.Marshal(&EvmRpcRequest{
		JsonRpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{block, true},
		Id:      "1",
	})

	response, error := http.Post(nodeUrl, "application/json", bytes.NewBuffer(params))

	if error != nil {
		fmt.Printf("%s %s\n", logHead, error.Error())
		return ""
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Printf("%s %s\n", logHead, response.Status)
		return ""
	} else {
		var data EvmGetBlockByNumberResponse

		decodeErr := json.NewDecoder(response.Body).Decode(&data)

		if decodeErr == nil {
			fmt.Printf("%s Block #%s have %d transactions\n", logHead, hexToDec(block), len(data.Result.Transactions))
			return data.Result
		}

		fmt.Printf("%s %s\n", logHead, decodeErr.Error())
	}

	return ""
}

func GetLatestBlock(nodeUrl string) string {
	const logHead = "[Celo Node RPC] GetLatestBlock -"

	params, _ := json.Marshal(&EvmRpcRequest{
		JsonRpc: "2.0",
		Method:  "eth_blockNumber",
		Params:  []interface{}{},
		Id:      "1",
	})

	response, error := http.Post(nodeUrl, "application/json", bytes.NewBuffer(params))

	if error != nil {
		fmt.Printf("%s %s\n", logHead, error.Error())
		return ""
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Printf("%s %s\n", logHead, response.Status)
		return ""
	} else {
		var data EvmGetLatestBlockResponse

		decodeErr := json.NewDecoder(response.Body).Decode(&data)

		if decodeErr == nil {
			fmt.Printf("%s %s (%s)\n", logHead, data.Result, hexToDec(data.Result))
			return data.Result
		}

		fmt.Printf("%s %s\n", logHead, decodeErr.Error())
	}

	return ""
}

func hexToDec(input string) string {
	postProcessed := strings.Replace(input, "0x", "", -1)
	decimal, error := strconv.ParseUint(postProcessed, 16, 64)

	if error != nil {
		fmt.Println(error)
	}

	return strconv.FormatUint(decimal, 10)
}
