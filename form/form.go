package form

import (
	"log"

	"github.com/charmbracelet/huh"
)

var (
	Currency string
	TopList  string
)

func GetFormData() (Currency string, TopList string) {
	form := huh.NewForm(
		huh.NewGroup(
			// Ask the user for a preferred currency and top list number.
			huh.NewSelect[string]().
				Title("Choose your preferred currency").
				Options(
					huh.NewOption("USD", "usd"),
					huh.NewOption("EUR", "eur"),
					huh.NewOption("GBP", "gbp"),
					huh.NewOption("INR", "inr"),
					huh.NewOption("AUD", "aud"),
				).
				Value(&Currency), // store the chosen option in the "Currency" variable

			// Let the user select top list of cryptocurrencies.
			huh.NewSelect[string]().
				Title("Top List (eg: top 10 cryptocurrencies)").
				Options(
					huh.NewOption("Top 10", "10").Selected(true),
					huh.NewOption("Top 25", "25"),
					huh.NewOption("Top 50", "50"),
				).
				Value(&TopList), // store the chosen option in the "TopList" variable
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	return Currency, TopList
}
