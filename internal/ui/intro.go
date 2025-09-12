package ui

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func RunIntro() (bool, error) {
	var start bool

	h1Style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#00B3FF")).
		Padding(0, 1)

	h2Style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Padding(0, 1)

	welcomeMessage := h1Style.Render("Welcome to CLI Crypto Prices!")
	infoMessage := h2Style.Render("This application allows you to view the latest cryptocurrency prices in your terminal. You can choose your preferred currency and the number of top cryptocurrencies to display.")

	fmt.Println()
	fmt.Println(welcomeMessage)
	fmt.Println()
	fmt.Println(infoMessage)
	fmt.Println()

	confirm := huh.NewConfirm().
		Title("Do you want to start?").
		Affirmative("Yes!").
		Negative("No.").
		Value(&start)

	err := confirm.Run()
	if err != nil {
		return false, err
	}

	return start, nil
}
