package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tranvictor/ethutils"
	"github.com/tranvictor/ethutils/txanalyzer"
)

var KYBER_ADDRESSES = []string{
	"0x964F35fAe36d75B1e72770e244F6595B68508CF5", // kyber v1
	"0x818E6FECD516Ecc3849DAf6845e3EC868087B755", // kyber v2
}

func isKyber(address string) bool {
	for _, a := range KYBER_ADDRESSES {
		if strings.ToLower(address) == strings.ToLower(a) {
			return true
		}
	}
	return false
}

func isTokenContract(address string) bool {
	tokendb, err := GetKyberTokenDB()
	if err != nil {
		fmt.Printf("Getting token db failed: %s\n", err)
		return false
	}
	return tokendb.IsToken(address)
}

func GetTokenDecimal(address string) (int, error) {
	tokendb, err := GetKyberTokenDB()
	if err != nil {
		fmt.Printf("Getting token db failed: %s\n", err)
		return 0, fmt.Errorf("couldn't get token db")
	}
	t, err := tokendb.GetToken(address)
	if err != nil {
		return 0, err
	}
	return t.Decimals, nil
}

func EventsFromInternals(tx *ethutils.TxInfo, wallet string) []Event {
	result := []Event{}
	for _, tx := range tx.InternalTxs {
		abig, ok := big.NewInt(0).SetString(tx.Value, 10)
		if !ok {
			fmt.Printf("Converting internal value to big int failed\n")
			continue
		}
		afloat := ethutils.BigToFloat(abig, 18)
		if afloat == 0 {
			continue
		}
		if l(tx.From) == l(wallet) {
			if l(tx.To) == l(wallet) {
				result = append(result, Event{
					Type:      SELF,
					InAsset:   "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
					InAmount:  afloat,
					OutAsset:  "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
					OutAmount: afloat,
				})
			} else {
				result = append(result, Event{
					Type:      WITHDRAW,
					OutAsset:  "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
					OutAmount: afloat,
				})
			}
		} else {
			if l(tx.To) == l(wallet) {
				result = append(result, Event{
					Type:     DEPOSIT,
					InAsset:  "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
					InAmount: afloat,
				})
			} else {
				result = append(result, Event{
					Type: UNKNOWN,
				})
			}
		}
	}
	return result
}

func EventFromKyberTrade(tx *ethutils.TxInfo, analyzedResult *txanalyzer.TxResult, wallet string) ([]Event, error) {
	result := []Event{}
	for _, log := range analyzedResult.Logs {
		// TODO: check log original address
		if log.Name == "ExecuteTrade" {
			outAsset := strings.Split(log.Data[0].Value, " ")[0]
			outDecimal, err := GetTokenDecimal(outAsset)
			if err != nil {
				return result, err
			}
			outAmount := ethutils.StringToFloat(strings.Split(log.Data[2].Value, " ")[0], int64(outDecimal))
			inAsset := strings.Split(log.Data[1].Value, " ")[0]
			inDecimal, err := GetTokenDecimal(inAsset)
			if err != nil {
				return result, err
			}
			inAmount := ethutils.StringToFloat(strings.Split(log.Data[3].Value, " ")[0], int64(inDecimal))

			result = append(result, Event{
				Type:      TRADE,
				InAsset:   inAsset,
				InAmount:  inAmount,
				OutAsset:  outAsset,
				OutAmount: outAmount,
			})
		}
	}
	return result, nil
}

func EventFromTrade(tx *ethutils.TxInfo, analyzedResult *txanalyzer.TxResult, wallet string) ([]Event, error) {
	return EventFromKyberTrade(tx, analyzedResult, wallet)
}

func EventsFromLogs(tx *ethutils.TxInfo, analyzedResult *txanalyzer.TxResult, wallet string) ([]Event, error) {
	result := []Event{}
	for _, log := range tx.Receipt.Logs {
		tokenAddr := log.Address.Hex()
		if isTokenContract(tokenAddr) && l(log.Topics[0].Hex()) == "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" {
			from := common.BigToAddress(log.Topics[1].Big()).Hex()
			to := common.BigToAddress(log.Topics[2].Big()).Hex()
			amount := big.NewInt(0).SetBytes(log.Data)
			decimal, err := GetTokenDecimal(tokenAddr)
			if err != nil {
				return nil, err
			}
			a := ethutils.BigToFloat(amount, int64(decimal))
			if l(from) == l(wallet) {
				if l(to) == l(wallet) {
					result = append(result, Event{
						Type:      SELF,
						InAsset:   tokenAddr,
						InAmount:  a,
						OutAsset:  tokenAddr,
						OutAmount: a,
					})
				} else {
					result = append(result, Event{
						Type:      WITHDRAW,
						OutAsset:  tokenAddr,
						OutAmount: a,
					})
				}
			} else {
				if l(to) == l(wallet) {
					result = append(result, Event{
						Type:     DEPOSIT,
						InAsset:  tokenAddr,
						InAmount: a,
					})
				} else {
					// ignored
				}
			}
		}
	}
	return result, nil
}

// func GetTokenTransferParams(tx *ethutils.TxInfo, result *txanalyzer.TxResult, decimal int64) (string, float64) {
// 	to, value := result.Params[0], result.Params[1]
// 	vHexWrapped := strings.Split(value.Value, " ")[1]
// 	vHex := vHexWrapped[1:(len(vHexWrapped) - 1)]
// 	vBig := ethutils.HexToBig(vHex)
// 	return strings.Split(to.Value, " ")[0], ethutils.BigToFloat(vBig, decimal)
// }

// func TypeOf(tx *ethutils.TxInfo, result *txanalyzer.TxResult, wallet string) int {
// 	if isKyber(tx.Tx.To().Hex()) {
// 		return TRADE
// 	} else {
// 		if isTokenContract(tx.Tx.To().Hex()) {
// 			if result.Method == "transfer" {
// 				// if the tx is erc20 transfer tx
// 				decimals, err := GetTokenDecimal(tx.Tx.To().Hex())
// 				if err != nil {
// 					fmt.Printf("Couldn't get token decimal: %s", err)
// 					return UNKNOWN
// 				}
// 				to, _ := GetTokenTransferParams(tx, result, int64(decimals))
// 				if strings.ToLower(tx.Tx.Extra.From.Hex()) == strings.ToLower(wallet) {
// 					if strings.ToLower(to) == strings.ToLower(wallet) {
// 						return SELF
// 					} else {
// 						return TOKEN_WITHDRAW
// 					}
// 				} else {
// 					if strings.ToLower(to) == strings.ToLower(wallet) {
// 						return TOKEN_DEPOSIT
// 					} else {
// 						return UNKNOWN
// 					}
// 				}
// 			} else {
// 				return UNKNOWN
// 			}
// 		} else {
// 			if tx.Tx.Value().Cmp(big.NewInt(0)) == 0 {
// 				// tx doesn't contain eth transfer and
// 				// doesn't call token contract
// 				return UNKNOWN
// 			} else {
// 				if strings.ToLower(tx.Tx.Extra.From.Hex()) == strings.ToLower(wallet) {
// 					if strings.ToLower(tx.Tx.To().Hex()) == strings.ToLower(wallet) {
// 						return SELF
// 					} else {
// 						return ETH_WITHDRAW
// 					}
// 				} else {
// 					if strings.ToLower(tx.Tx.To().Hex()) == strings.ToLower(wallet) {
// 						return ETH_DEPOSIT
// 					} else {
// 						return UNKNOWN
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return UNKNOWN
// }
