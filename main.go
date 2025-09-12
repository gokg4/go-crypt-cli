package main

import (
	"cli-view-crypto-prices/internal/ui"
	"cli-view-crypto-prices/internal/utils"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	currencyFlag := flag.String("currency", "", "The currency to display prices in (e.g., 'usd', 'eur', 'jpy')")
	limitFlag := flag.Int("limit", 0, "The number of cryptocurrencies to display")
	flag.Parse()

	useFlags := *currencyFlag != "" && *limitFlag > 0

	if !useFlags {
		start, err := ui.RunIntro()
		if err != nil {
			log.Fatal(err)
		}
		if !start {
			os.Exit(0)
		}
	}

	var currency string
	var limit int

	if useFlags {
		currency = *currencyFlag
		limit = *limitFlag
	}

	for {
		var err error
		if !useFlags {
			var limitStr string
			currency, limitStr, err = ui.RunForm()
			if err != nil {
				log.Fatal(err)
			}

			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				log.Fatal(err)
			}
		}

		m := ui.NewModel(currency, limit)
		p := tea.NewProgram(m)

		finalModel, err := p.Run()
		if err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

		tableModel := finalModel.(*ui.Model)
		if !tableModel.ShowForm {
			break
		}

		useFlags = false
		utils.ClearScreen()
	}
}
