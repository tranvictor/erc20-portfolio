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

func GetAllTxInfo(hashes []string, internalTxs map[string][]ethutils.InternalTx) ([]*ethutils.TxInfo, error) {
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
				if len(txinfo.InternalTxs) == 0 && len(internalTxs[l(hash)]) > 0 {
					txinfo.InternalTxs = internalTxs[l(hash)]
					err = db.StoreTxs([]*ethutils.TxInfo{txinfo})
					if err != nil {
						return []*ethutils.TxInfo{}, err
					}
				}
				result = append(result, txinfo)
			} else {
				// get receipt from blockchain
				fmt.Printf("%d. getting receipt for %s\n", i, hash)
				tx, err := rd.TxInfoFromHash(hash)
				if err != nil {
					return []*ethutils.TxInfo{}, err
				}
				if len(tx.InternalTxs) == 0 && len(internalTxs[l(hash)]) > 0 {
					tx.InternalTxs = internalTxs[l(hash)]
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
