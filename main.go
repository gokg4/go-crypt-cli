package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	marketdata "gocryptocli/marketData"
	createtable "gocryptocli/table"
)

const VsCurrencyDefault = "inr"

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("What is your preferred currency? eg: usd,eur,gbp,aud,inr...etc")
	fmt.Printf("> ")
	scanner.Scan()
	currency := strings.ToLower(scanner.Text())

	var data marketdata.GeckoMarketData
	if currency == "" {
		currency = VsCurrencyDefault
	}
	data = marketdata.MarketData(currency)

	createtable.CreateTable(data, currency)
}
