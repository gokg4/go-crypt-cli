package main

import (
	"gocryptocli/form"
	marketdata "gocryptocli/marketData"
	createtable "gocryptocli/table"
)

const VsCurrencyDefault = "inr"

func main() {

	currency, topList := form.GetFormData()

	var data marketdata.GeckoMarketData
	if currency == "" {
		currency = VsCurrencyDefault
	}
	data = marketdata.MarketData(currency, topList)

	createtable.CreateTable(data, currency)
}
