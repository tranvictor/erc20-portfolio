package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/tranvictor/ethutils"
)

const (
	TX_DB_FILE string = "tx_db.json"
)

type TxData struct {
	Hash        string `json:"hash"`
	BlockNumber string `json:"blockNumber"`
}

type TxDataJSONDB struct {
	Data map[string]*ethutils.TxInfo `json:"data"`
}

func (self *TxDataJSONDB) StoreTxs(txs []*ethutils.TxInfo) error {
	for _, tx := range txs {
		hash := strings.ToLower(tx.Receipt.TxHash.Hex())
		self.Data[hash] = tx
	}
	return self.Persist()
}

func (self *TxDataJSONDB) GetTx(hash string) (*ethutils.TxInfo, error) {
	r, _ := self.Data[strings.ToLower(hash)]
	return r, nil
}

func (self *TxDataJSONDB) Persist() error {
	jsonData, err := json.MarshalIndent(self, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(TX_DB_FILE, jsonData, 0644)
}

func NewTxDataJSONDB() (*TxDataJSONDB, error) {
	jsonData, err := ioutil.ReadFile(TX_DB_FILE)
	if err != nil {
		return nil, err
	}
	db := &TxDataJSONDB{
		Data: map[string]*ethutils.TxInfo{},
	}
	err = json.Unmarshal(jsonData, db)
	if err != nil {
		return nil, err
	}
	return db, nil
}
