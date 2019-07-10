package main

func main() {
	wallet := "0x5ed3707FF33a3DFC71f6fa109Fa6eF7D9B5DAC69"
	// wallet := "0xA843c8ef4F35aD775a8169a0B0Aebb048a7A7572"
	// wallet := "0xf214dde57f32f3f34492ba3148641693058d4a9e"
	// wallet := "0x06173b3384c75ead92fceb9297e2593a33fdbff2"
	// wallet := "0xAEAFE589E6933B0D1B8fA5D1941813472F68C308"
	// wallet := "0x3D5c8B827E40cFE0b3406Ae5EF57c6828DA02844"
	// wallet := "0xa6c883e2dde82fbed20e025bd717a6b7f34f5e6e"

	portfolio, err := NewPortfolioFromFile(wallet)
	if err != nil {
		panic(err)
	}
	portfolio.Update()
	portfolio.Print()
}
