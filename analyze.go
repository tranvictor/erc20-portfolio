package main

import (
	"fmt"
)

type PortfolioResult struct {
	Fee               float64            `json:"total_fee"`
	currentInvestment map[string]float64 `json:"currentInvestment"`
	currentPortfolio  map[string]float64 `json:"currentPortfolio"`
	withdrawals       map[string]float64 `json:"withdrawals"`
}

func NewPortfolioResult() *PortfolioResult {
	return &PortfolioResult{
		Fee:               0,
		currentInvestment: map[string]float64{},
		currentPortfolio:  map[string]float64{},
		withdrawals:       map[string]float64{},
	}
}

func (self *PortfolioResult) RegisterFee(fee float64) {
	self.Fee += fee
}

func (self *PortfolioResult) registerTrade(e Event) {
	tin, err := e.InToken()
	if err != nil {
		fmt.Printf("getting in token failed: %s\n", err)
		return
	}
	tout, err := e.OutToken()
	if err != nil {
		fmt.Printf("getting in token failed: %s\n", err)
		return
	}

	_, found := self.currentPortfolio[tout.Symbol]
	if found {
		self.currentPortfolio[tout.Symbol] += -e.OutAmount
	} else {
		self.currentPortfolio[tout.Symbol] = -e.OutAmount
	}
	_, found = self.currentPortfolio[tin.Symbol]
	if found {
		self.currentPortfolio[tin.Symbol] += e.InAmount
	} else {
		self.currentPortfolio[tin.Symbol] = e.InAmount
	}
}

func (self *PortfolioResult) registerWithdraw(e Event) {
	tout, err := e.OutToken()
	if err != nil {
		fmt.Printf("getting in token failed: %s\n", err)
		return
	}

	_, found := self.currentPortfolio[tout.Symbol]
	if found {
		self.currentPortfolio[tout.Symbol] += -e.OutAmount
	} else {
		self.currentPortfolio[tout.Symbol] = -e.OutAmount
	}
	_, found = self.withdrawals[tout.Symbol]
	if found {
		self.withdrawals[tout.Symbol] += e.OutAmount
	} else {
		self.withdrawals[tout.Symbol] = e.OutAmount
	}
}

func (self *PortfolioResult) registerDeposit(e Event) {
	tin, err := e.InToken()
	if err != nil {
		fmt.Printf("getting in token failed: %s\n", err)
		return
	}
	_, found := self.currentInvestment[tin.Symbol]
	if found {
		self.currentInvestment[tin.Symbol] += e.InAmount
	} else {
		self.currentInvestment[tin.Symbol] = e.InAmount
	}
	_, found = self.currentPortfolio[tin.Symbol]
	if found {
		self.currentPortfolio[tin.Symbol] += e.InAmount
	} else {
		self.currentPortfolio[tin.Symbol] = e.InAmount
	}
}

func (self *PortfolioResult) RegisterEvent(e Event) {
	switch e.Type {
	case DEPOSIT:
		self.registerDeposit(e)
	case WITHDRAW:
		self.registerWithdraw(e)
	case TRADE:
		self.registerTrade(e)
	case SELF, UNKNOWN:
	}
}

func (self *PortfolioResult) Withdrew() map[string]float64 {
	return self.withdrawals
}

func (self *PortfolioResult) Investment() map[string]float64 {
	return self.currentInvestment
}

func (self *PortfolioResult) Portfolio() map[string]float64 {
	return self.currentPortfolio
}

func (self *PortfolioResult) Pnl() float64 {
	return 0
}

func (self *PortfolioResult) PnlUSD() float64 {
	return 0
}

func (self *PortfolioResult) TotalFee() float64 {
	return self.Fee
}
