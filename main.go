package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/nighostchris/tx-track-engine-go/database"
)

func main() {
	connectionParams := database.DatabaseConnectionParams{
		Username: "root",
		Password: "root",
		Host:     "0.0.0.0",
		Port:     "5432",
		Database: "postgres",
	}

	db, err := database.Connect(connectionParams)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(db)
	}

	// const nodeUrl = "https://alfajores-forno.celo-testnet.org"

	// threads := make(chan bool, 5)
	// var lastProcessedBlock = "0x8bc540"

	// for {
	// 	blockNumber := GetLatestBlock(nodeUrl)

	// 	if blockNumber >= lastProcessedBlock && len(threads) < 5 {
	// 		var poolCapacity = 5 - len(threads)

	// 		for i := 0; i < poolCapacity; i++ {
	// 			midwayBlock, _ := strconv.ParseUint(lastProcessedBlock[2:], 16, 64)
	// 			targetBlock := "0x" + strconv.FormatUint(midwayBlock+1, 16)

	// 			go GetBlockByNumber(nodeUrl, targetBlock, threads)

	// 			threads <- true
	// 			lastProcessedBlock = targetBlock
	// 		}
	// 	}

	// 	time.Sleep(5 * time.Second)
	// }
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

func GetBlockByNumber(nodeUrl string, block string, pool chan bool) interface{} {
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
		<-pool
		return ""
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Printf("%s %s\n", logHead, response.Status)
		<-pool
		return ""
	} else {
		var data EvmGetBlockByNumberResponse

		decodeErr := json.NewDecoder(response.Body).Decode(&data)

		if decodeErr == nil {
			fmt.Printf("%s Block #%s have %d transactions\n", logHead, hexToDec(block), len(data.Result.Transactions))
			<-pool
			return data.Result
		}

		fmt.Printf("%s %s\n", logHead, decodeErr.Error())
	}

	<-pool
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
