package main

import (
	"fmt"
	"gocryptocli/form"
	marketData "gocryptocli/marketData"
	"gocryptocli/table"
	"time"

	"github.com/charmbracelet/huh/spinner"
)

const VsCurrencyDefault = "inr"

func main() {

	currency, topList := form.GetFormData()

	var data marketData.GeckoMarketData
	if currency == "" {
		currency = VsCurrencyDefault
	}

	action := func() {
		data = marketData.GetMarketData(currency, topList)
		time.Sleep(2 * time.Second)
	}

	fmt.Println()
	_ = spinner.New().Title("Getting data...").Action(action).Run()

	table.CreateTable(data, currency)
}
