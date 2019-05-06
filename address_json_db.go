package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/tranvictor/ethutils/reader"
)

const (
	ADDRESS_DB_FILE string = "address_db.json"
)

var AddressDB *AddressJSONDB
var adbfOnce sync.Once

type AddressData struct {
	IsContract bool   `json:"isContract"`
	ABI        string `json:"abi"`
}

type AddressJSONDB struct {
	Data map[string]*AddressData `json:"data"`
}

func (self *AddressJSONDB) GetAddress(address string) (bool, string, error) {
	data, found := self.Data[strings.ToLower(address)]
	if !found {
		fmt.Printf("Looking up address data for %s\n", address)
		rd := reader.NewEthReader()
		code, err := rd.GetCode(address)
		if err != nil {
			return false, "", err
		}
		isContract := len(code) > 0
		var abi string
		if isContract {
			abi, err = rd.GetABIString(address)
			if err != nil {
				return false, "", err
			}
		}
		self.Data[strings.ToLower(address)] = &AddressData{
			IsContract: isContract,
			ABI:        abi,
		}
		defer self.Persist()
		return isContract, abi, nil
	} else {
		return data.IsContract, data.ABI, nil
	}
}

func (self *AddressJSONDB) Persist() error {
	jsonData, err := json.MarshalIndent(self, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ADDRESS_DB_FILE, jsonData, 0644)
}

func NewAddressJSONDB() (*AddressJSONDB, error) {
	var err error
	adbfOnce.Do(func() {
		var jsonData []byte
		jsonData, err = ioutil.ReadFile(ADDRESS_DB_FILE)
		if err != nil {
			return
		}
		db := &AddressJSONDB{
			Data: map[string]*AddressData{},
		}
		err = json.Unmarshal(jsonData, db)
		if err != nil {
			return
		}
		AddressDB = db
	})
	return AddressDB, err
}
