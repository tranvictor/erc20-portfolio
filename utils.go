package main

import (
	"fmt"
	"sort"
	"strings"
)

func l(s string) string { return strings.ToLower(s) }

func PrintBalances(balances map[string]float64) {
	assets := []string{}
	for a, _ := range balances {
		assets = append(assets, a)
	}
	sort.Strings(assets)
	for _, a := range assets {
		fmt.Printf("%s: %f\n", a, balances[a])
	}
}
