package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

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

func GetAllTxsFromEtherscan(wallet string) ([]*TxData, error) {
	normalTxs, err := GetAllNormalTxsFromEtherscan(wallet)
	if err != nil {
		return []*TxData{}, err
	}
	tokenTxs, err := GetAllTokenTxsFromEtherscan(wallet)
	if err != nil {
		return []*TxData{}, err
	}
	temp := map[string]*TxData{}
	for _, tx := range normalTxs {
		temp[strings.ToLower(tx.Hash)] = tx
	}
	for _, tx := range tokenTxs {
		temp[strings.ToLower(tx.Hash)] = tx
	}
	result := []*TxData{}
	for _, tx := range temp {
		result = append(result, tx)
	}
	sort.Sort(ListOfTxs(result))
	return result, nil
}
