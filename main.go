package main

import "fmt"

func PrintPortfolio(portfolio map[string]float64) {
	for token, balance := range portfolio {
		if balance >= 0.001 {
			fmt.Printf("    %s: %f\n", token, balance)
		}
	}
}

func main() {
	wallet := "0x5ed3707FF33a3DFC71f6fa109Fa6eF7D9B5DAC69"
	// wallet := "0xA843c8ef4F35aD775a8169a0B0Aebb048a7A7572"

	portfolio, err := NewPortfolioFromFile(wallet)
	if err != nil {
		panic(err)
	}
	portfolio.Update()
	portfolio.Print()
}
