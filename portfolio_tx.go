package main

type PortfolioTx struct {
	Hash   string  `json:"hash"`
	Events []Event `json:"events"`
	Fee    float64 `json:"fee"`
}
