package main

type EventType int

const (
	UNKNOWN EventType = iota
	DEPOSIT
	TRADE
	WITHDRAW
	SELF
)

type Event struct {
	Type      EventType `json:"type"`
	InAsset   string    `json:"in_asset"`
	InAmount  float64   `json:"in_amount"`
	OutAsset  string    `json:"out_asset"`
	OutAmount float64   `json:"out_amount"`
}

func (self *Event) InToken() (*Token, error) {
	tokendb, err := GetKyberTokenDB()
	if err != nil {
		return nil, err
	}
	return tokendb.GetToken(self.InAsset)
}

func (self *Event) OutToken() (*Token, error) {
	tokendb, err := GetKyberTokenDB()
	if err != nil {
		return nil, err
	}
	return tokendb.GetToken(self.OutAsset)
}
