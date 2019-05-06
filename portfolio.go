package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/tranvictor/ethutils"
	"github.com/tranvictor/ethutils/txanalyzer"
)

type Portfolio struct {
	Address string        `json:"address"`
	Txs     []PortfolioTx `json:"txs"`
	txmap   map[string]bool
}

func CreatePortfolio(address string) (*Portfolio, error) {
	result := &Portfolio{
		Address: address,
		Txs:     []PortfolioTx{},
		txmap:   map[string]bool{},
	}
	err := result.Update()
	if err != nil {
		return nil, err
	}
	err = result.Persist()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func NewPortfolioFromFile(address string) (*Portfolio, error) {
	jsonData, err := ioutil.ReadFile(l(address))
	if err != nil {
		result, err := CreatePortfolio(address)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	result := &Portfolio{
		Address: address,
		Txs:     []PortfolioTx{},
		txmap:   map[string]bool{},
	}
	err = json.Unmarshal(jsonData, result)
	if err != nil {
		return nil, err
	}
	for _, tx := range result.Txs {
		result.txmap[l(tx.Hash)] = true
	}
	return result, nil
}

func (self *Portfolio) Update() error {
	txs, err := GetAllTxsFromEtherscan(self.Address)
	if err != nil {
		return err
	}

	hashes := []string{}
	for _, tx := range txs {
		hashes = append(hashes, tx.Hash)
	}
	allTxInfo, err := GetAllTxInfo(hashes)
	if err != nil {
		return err
	}
	err = self.extractEvents(allTxInfo)
	if err != nil {
		return err
	}
	return self.Persist()
}

func (self *Portfolio) extractEvents(txs []*ethutils.TxInfo) error {
	for _, tx := range txs {
		if !self.txmap[l(tx.Tx.Hash().Hex())] {
			err := self.extractEventsInOneTx(tx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (self *Portfolio) extractEventsInOneTx(tx *ethutils.TxInfo) error {
	result := PortfolioTx{
		Hash:   tx.Tx.Hash().Hex(),
		Events: []Event{},
	}
	if l(tx.Tx.Extra.From.Hex()) == l(self.Address) {
		result.Fee = ethutils.BigToFloat(tx.GasCost(), 18)
	}
	var err error
	addressDB, err := NewAddressJSONDB()
	if err != nil {
		return err
	}
	isContract, abiStr, err := addressDB.GetAddress(tx.Tx.To().Hex())
	if err != nil {
		return err
	}
	var a abi.ABI
	if isContract && abiStr != "Contract source code not verified" {
		a, err = abi.JSON(strings.NewReader(abiStr))
		if err != nil {
			return err
		}
	}
	analyzer := txanalyzer.NewAnalyzer()
	analyzedResult := analyzer.AnalyzeOffline(tx, &a, isContract)
	if isKyber(tx.Tx.To().Hex()) {
		// TODO: parse kyber trade
		result.Events = append(result.Events, Event{
			Type:      TRADE,
			InAsset:   "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
			InAmount:  0.0,
			OutAsset:  "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
			OutAmount: 0.0,
		})
	} else {
		if isContract {
			if tx.Tx.Value().Cmp(big.NewInt(0)) != 0 {
				if l(tx.Tx.Extra.From.Hex()) == l(self.Address) {
					if l(tx.Tx.To().Hex()) == l(self.Address) {
						result.Events = append(result.Events, Event{
							Type: SELF,
						})
					} else {
						result.Events = append(result.Events, Event{
							Type:      WITHDRAW,
							OutAsset:  "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
							OutAmount: ethutils.BigToFloat(tx.Tx.Value(), 18),
						})
					}
				} else {
					if l(tx.Tx.To().Hex()) == l(self.Address) {
						result.Events = append(result.Events, Event{
							Type:     DEPOSIT,
							InAsset:  "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
							InAmount: ethutils.BigToFloat(tx.Tx.Value(), 18),
						})
					} else {
						result.Events = append(result.Events, Event{
							Type: UNKNOWN,
						})
					}
				}
			}
			events, err := EventsFromLogs(tx, analyzedResult, self.Address)
			if err != nil {
				return err
			} else {
				result.Events = append(result.Events, events...)
			}
			// TODO: result.Events = append(result.Events, EventsFromInternals(tx, self.Address)...)
		} else {
			if tx.Tx.Value().Cmp(big.NewInt(0)) == 0 {
				result.Events = append(result.Events, Event{
					Type: UNKNOWN,
				})
			} else {
				if l(tx.Tx.Extra.From.Hex()) == l(self.Address) {
					if l(tx.Tx.To().Hex()) == l(self.Address) {
						result.Events = append(result.Events, Event{
							Type: SELF,
						})
					} else {
						result.Events = append(result.Events, Event{
							Type:      WITHDRAW,
							OutAsset:  "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
							OutAmount: ethutils.BigToFloat(tx.Tx.Value(), 18),
						})
					}
				} else {
					if l(tx.Tx.To().Hex()) == l(self.Address) {
						result.Events = append(result.Events, Event{
							Type:     DEPOSIT,
							InAsset:  "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
							InAmount: ethutils.BigToFloat(tx.Tx.Value(), 18),
						})
					} else {
						result.Events = append(result.Events, Event{
							Type: UNKNOWN,
						})
					}
				}
			}
		}
	}
	self.Txs = append(self.Txs, result)
	self.txmap[l(tx.Tx.Hash().Hex())] = true
	return nil
}

func (self *Portfolio) analyze() *PortfolioResult {
	result := NewPortfolioResult()
	return result
}

func (self *Portfolio) Print() {
	result := self.analyze()

	for i, tx := range self.Txs {
		fmt.Printf("%d. %s - (%f ETH fee)\n", i, tx.Hash, tx.Fee)
		for j, e := range tx.Events {
			switch e.Type {
			case UNKNOWN:
				fmt.Printf("-- %d. UNKNOWN event\n", j)
			case DEPOSIT:
				t, err := e.InToken()
				if err != nil {
					fmt.Printf("-- %d. Getting token info failed: %s\n", j, err)
				} else {
					fmt.Printf("-- %d. DEPOSIT %f %s\n", j, e.InAmount, t.Symbol)
				}
			case WITHDRAW:
				t, err := e.OutToken()
				if err != nil {
					fmt.Printf("-- %d. Getting token info failed: %s\n", j, err)
				} else {
					fmt.Printf("-- %d. WITHDRAW %f %s\n", j, e.OutAmount, t.Symbol)
				}
			case TRADE:
				fmt.Printf("-- %d. TRADE event\n", j)
			case SELF:
				fmt.Printf("-- %d. SELF event\n", j)
			}
		}
	}

	fmt.Printf("Portfolio summary:\n")
	fmt.Printf("1. Withdrew portfolio:\n")
	PrintPortfolio(result.Withdrew())
	fmt.Printf("2. Init portfolio:\n")
	PrintPortfolio(result.Starting())
	fmt.Printf("3. Current portfolio:\n")
	PrintPortfolio(result.Current())
	fmt.Printf("4. Pnl in ETH: %f ETH\n", result.Pnl())
	fmt.Printf("5. Pnl in USD: %f USD\n", result.PnlUSD())
	fmt.Printf("6. Fee expense: %f ETH\n", result.TotalFee())
}

func (self *Portfolio) Persist() error {
	return nil
	// jsonData, err := json.MarshalIndent(self, "", "  ")
	// if err != nil {
	// 	return err
	// }
	// return ioutil.WriteFile(self.Address, jsonData, 0644)
}
