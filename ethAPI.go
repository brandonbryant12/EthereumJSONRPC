package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	//"time"
	//"reflect"
)

type Payment struct {
	Currency string
	Address  string
	Amount   string
	Hash     string
}

type Payload struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type Block struct {
	Difficulty       string        `json:"difficulty"`
	ExtraData        string        `json:"extraData"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Hash             string        `json:"hash"`
	LogsBloom        string        `json:"logsBloom"`
	Miner            string        `json:"miner"`
	MixHash          string        `json:"mixHash"`
	Nonce            string        `json:"nonce"`
	Number           string        `json:"number"`
	ParentHash       string        `json:"parentHash"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	SHA3Uncles       string        `json:"sha3Uncles"`
	Size             string        `json:"size"`
	StateRoot        string        `json:"stateRoot"`
	Timestamp        string        `json:"timeStamp"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	Transactions     []Transaction `json:"transactions"`
	TransactionsRoot string        `json:"transactions_root"`
	Uncles           []string      `json:"uncles"`
}

type Transaction struct {
	BlockHash        string  `json:"blockHash"`
	BlockNumber      *string `json:"blockNumber"`
	Gas              string  `json:"gas"`
	GasPrice         string  `json:"gasPrice"`
	Hash             string  `json:"hash"`
	Input            string  `json:"input"`
	Nonce            string  `json:"nonce"`
	To               string  `json:"to"`
	R                string  `json:"r"`
	S                string  `json:"s"`
	TransactionIndex string  `json:"transactionIndex"`
	V                string  `json:"v"`
	Value            string  `json:"value"`
}

func (tx *Transaction) String() string {
	return fmt.Sprintf("to: %v\ngas: %v\ngasPrice: %v\nvalue: %v\nblockHash: %v\nblockNumber: %v\nhash: %v\ninput: %v\nnounce: %v\nr: %v\ns: %v\nv: %v\ntransactionIndex: %v\n",
		tx.To,
		tx.Gas,
		tx.GasPrice,
		tx.Value,
		tx.BlockHash,
		tx.BlockNumber,
		tx.Hash,
		tx.Input,
		tx.Nonce,
		tx.R,
		tx.S,
		tx.V,
		tx.TransactionIndex,
	)
}

//Refactor to make this method take an array of any type and size.
func setParams(blockNumber string, verbose bool) []interface{} {
	params := make([]interface{}, 2)
	params[0] = blockNumber
	params[1] = verbose
	return params
}

func startPolling() {

}

func handleRequest(req *http.Request) string {

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	
	s := strings.SplitAfter(string(content), `"result":`)[1]
	return s[:len(s)-len("}")]
}

//Takes and array of Transaction objects and outputs an array of Payment structs
func processTxs(txs []Transaction) []Payment {
	var payments []Payment
	for i := range txs {
		payments = append(payments, Payment{
							Currency: "ETH",
							Address: txs[i].To,
							Amount:	txs[i].Value,
							Hash:	txs[i].Hash})
	}
	return payments
}

func processBlock(rawResponse string) Block {
	var block Block
        json.Unmarshal([]byte(rawResponse), &block)
	return block
}

func main() {

	params := setParams("latest", true)

	//Refactor to generate a payload given any method
	data := Payload{Jsonrpc: "2.0", Method: "eth_getBlockByNumber", Params: params, ID: 1}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}

	body := bytes.NewReader(payloadBytes)
	url := "https://mainnet.infura.io/v3/924c0f97172441a28a5b5270db968474"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/json")
	result := handleRequest(req)
	block := processBlock(result)
/*
	payments := processTxs(block.Transactions)
	for i := range payments {
		fmt.Println(payments[i])
	}
*/	

}
