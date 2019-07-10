package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/tranvictor/ethutils"
)

type internaltx struct {
	Hash        string `json:"hash"`
	BlockNumber string `json:"blockNumber"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
}

type EtherscanInternalTxResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []*internaltx `json:"result"`
}

type EtherscanResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Result  []*TxData `json:"result"`
}

func GetAllTokenTxsFromEtherscan(wallet string) ([]*TxData, error) {
	url := fmt.Sprintf(
		"http://api.etherscan.io/api?module=account&action=tokentx&address=%s&startblock=%d&endblock=%d&sort=desc&apikey=%s&isError=0",
		wallet,
		0,
		9999999999,
		"DS24ZMNDG8CGKNU1QAU32WFS4DZKJ3PE3J",
	)
	resp, err := http.Get(url)
	if err != nil {
		return []*TxData{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []*TxData{}, err
	}
	txs := EtherscanResponse{}
	err = json.Unmarshal(body, &txs)
	if err != nil {
		return []*TxData{}, err
	}
	return txs.Result, nil
}

func GetAllInternalTxsFromEtherscan(wallet string) ([]*internaltx, error) {
	url := fmt.Sprintf(
		"http://api.etherscan.io/api?module=account&action=txlistinternal&address=%s&startblock=%d&endblock=%d&sort=desc&apikey=%s&isError=0",
		wallet,
		0,
		9999999999,
		"DS24ZMNDG8CGKNU1QAU32WFS4DZKJ3PE3J",
	)
	resp, err := http.Get(url)
	if err != nil {
		return []*internaltx{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []*internaltx{}, err
	}
	txs := EtherscanInternalTxResponse{}
	err = json.Unmarshal(body, &txs)
	if err != nil {
		return []*internaltx{}, err
	}
	return txs.Result, nil
}

func GetAllNormalTxsFromEtherscan(wallet string) ([]*TxData, error) {
	url := fmt.Sprintf(
		"http://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=%d&endblock=%d&sort=desc&apikey=%s&isError=0",
		wallet,
		0,
		9999999999,
		"DS24ZMNDG8CGKNU1QAU32WFS4DZKJ3PE3J",
	)
	resp, err := http.Get(url)
	if err != nil {
		return []*TxData{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []*TxData{}, err
	}
	txs := EtherscanResponse{}
	err = json.Unmarshal(body, &txs)
	if err != nil {
		return []*TxData{}, err
	}
	return txs.Result, nil
}

type ListOfTxs []*TxData

func (self ListOfTxs) Len() int {
	return len(self)
}

func (self ListOfTxs) Less(i, j int) bool {
	itx := self[i]
	jtx := self[j]
	iblockno, _ := strconv.Atoi(itx.BlockNumber)
	jblockno, _ := strconv.Atoi(jtx.BlockNumber)
	return iblockno < jblockno
}

func (self ListOfTxs) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func GetAllTxsFromEtherscan(wallet string) ([]*TxData, map[string][]ethutils.InternalTx, error) {
	fmt.Printf("Getting all tx hashes of this wallet from etherscan...\n")
	normalTxs, err := GetAllNormalTxsFromEtherscan(wallet)
	if err != nil {
		return []*TxData{}, map[string][]ethutils.InternalTx{}, err
	}
	fmt.Printf("Getting all token txs related to this wallet from etherscan...\n")
	tokenTxs, err := GetAllTokenTxsFromEtherscan(wallet)
	if err != nil {
		return []*TxData{}, map[string][]ethutils.InternalTx{}, err
	}
	fmt.Printf("Getting all internal txs related to this wallet from etherscan...\n")
	internalTxs, err := GetAllInternalTxsFromEtherscan(wallet)
	if err != nil {
		return []*TxData{}, map[string][]ethutils.InternalTx{}, err
	}
	temp := map[string]*TxData{}
	for _, tx := range normalTxs {
		temp[strings.ToLower(tx.Hash)] = tx
	}
	for _, tx := range tokenTxs {
		temp[strings.ToLower(tx.Hash)] = tx
	}
	for _, tx := range internalTxs {
		temp[strings.ToLower(tx.Hash)] = &TxData{
			Hash:        tx.Hash,
			BlockNumber: tx.BlockNumber,
		}
	}
	result := []*TxData{}
	for _, tx := range temp {
		result = append(result, tx)
	}
	sort.Sort(ListOfTxs(result))
	internalTxResult := map[string][]ethutils.InternalTx{}
	for _, it := range internalTxs {
		if _, found := internalTxResult[l(it.Hash)]; found {
			internalTxResult[l(it.Hash)] = append(internalTxResult[l(it.Hash)], ethutils.InternalTx{
				From:  it.From,
				To:    it.To,
				Value: it.Value,
			})
		} else {
			internalTxResult[l(it.Hash)] = []ethutils.InternalTx{
				ethutils.InternalTx{
					From:  it.From,
					To:    it.To,
					Value: it.Value,
				},
			}
		}
	}
	return result, internalTxResult, nil
}
