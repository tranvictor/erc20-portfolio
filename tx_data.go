package main

import (
	"fmt"

	"github.com/tranvictor/ethutils"
	"github.com/tranvictor/ethutils/reader"
)

type TxDB interface {
	GetTx(hash string) (*ethutils.TxInfo, error)
	StoreTxs(txs []*ethutils.TxInfo) error
}

func GetTxDB() (TxDB, error) {
	return NewTxDataJSONDB()
}

func GetAllTxInfo(hashes []string) ([]*ethutils.TxInfo, error) {
	result := []*ethutils.TxInfo{}
	db, err := GetTxDB()
	rd := reader.NewEthReader()
	if err != nil {
		return []*ethutils.TxInfo{}, err
	}
	for i, hash := range hashes {
		txinfo, err := db.GetTx(hash)
		if err != nil {
			return []*ethutils.TxInfo{}, err
		} else {
			if txinfo != nil {
				result = append(result, txinfo)
			} else {
				// get receipt from blockchain
				fmt.Printf("%d. getting receipt for %s\n", i, hash)
				tx, err := rd.TxInfoFromHash(hash)
				if err != nil {
					return []*ethutils.TxInfo{}, err
				}
				err = db.StoreTxs([]*ethutils.TxInfo{&tx})
				if err != nil {
					return []*ethutils.TxInfo{}, err
				}
				result = append(result, &tx)
			}
		}
	}
	return result, nil
}
