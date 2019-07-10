package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

var KyberTokenDB *TokenDB
var once sync.Once

type Token struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Decimals int    `json:"decimals"`
}

type TokenDB struct {
	Data map[string]Token `json:"data"`
}

func (self *TokenDB) GetToken(address string) (*Token, error) {
	t, found := self.Data[l(address)]
	if !found {
		return nil, fmt.Errorf("token %s not found", address)
	}
	return &t, nil
}

func (self *TokenDB) IsToken(address string) bool {
	_, found := self.Data[l(address)]
	return found
}

func NewTokenDB() *TokenDB {
	return &TokenDB{
		Data: map[string]Token{},
	}
}

type kyberresp struct {
	Data  []Token `json:"data"`
	Error bool    `json:"error"`
}

func GetKyberTokenDB() (*TokenDB, error) {
	var err error
	once.Do(func() {
		url := "https://api.kyber.network/currencies?include_delisted=true"
		var resp *http.Response
		resp, err = http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		tokens := kyberresp{}
		err = json.Unmarshal(body, &tokens)
		if err != nil {
			return
		}
		KyberTokenDB = NewTokenDB()
		for _, token := range tokens.Data {
			KyberTokenDB.Data[l(token.Address)] = token
		}
	})
	return KyberTokenDB, err
}
