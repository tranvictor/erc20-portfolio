package main

type PortfolioResult struct {
}

func NewPortfolioResult() *PortfolioResult {
	return &PortfolioResult{}
}

func (self *PortfolioResult) Withdrew() map[string]float64 {
	return map[string]float64{}
}

func (self *PortfolioResult) Starting() map[string]float64 {
	return map[string]float64{}
}

func (self *PortfolioResult) Current() map[string]float64 {
	return map[string]float64{}
}

func (self *PortfolioResult) Pnl() float64 {
	return 0
}

func (self *PortfolioResult) PnlUSD() float64 {
	return 0
}

func (self *PortfolioResult) TotalFee() float64 {
	return 0
}
