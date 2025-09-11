package ui

import (
	"github.com/charmbracelet/huh"
)

func RunForm() (string, string, error) {
	var currency, limit string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your preferred currency").
				Options(
					huh.NewOption("USD", "usd"),
					huh.NewOption("EUR", "eur"),
					huh.NewOption("GBP", "gbp"),
					huh.NewOption("INR", "inr"),
					huh.NewOption("AUD", "aud"),
				).
				Value(&currency),

			huh.NewSelect[string]().
				Title("Top List (eg: top 10 cryptocurrencies)").
				Options(
					huh.NewOption("Top 10", "10").Selected(true),
					huh.NewOption("Top 25", "25"),
					huh.NewOption("Top 50", "50"),
				).
				Value(&limit),
		),
	)

	err := form.Run()
	if err != nil {
		return "", "", err
	}
	return currency, limit, nil
}
