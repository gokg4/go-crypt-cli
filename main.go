package main

import (
	"gocryptocli/form"
	marketdata "gocryptocli/marketData"
	createtable "gocryptocli/table"
)

const VsCurrencyDefault = "inr"

func main() {
	// scanner := bufio.NewScanner(os.Stdin)
	// fmt.Println("What is your preferred currency? eg: usd,eur,gbp,aud,inr...etc")
	// fmt.Printf("> ")
	// scanner.Scan()
	// currency := strings.ToLower(scanner.Text())

	currency, topList := form.GetFormData()

	var data marketdata.GeckoMarketData
	if currency == "" {
		currency = VsCurrencyDefault
	}
	data = marketdata.MarketData(currency, topList)

	createtable.CreateTable(data, currency)
}
