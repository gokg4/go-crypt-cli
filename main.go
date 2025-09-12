package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"cli-view-crypto-prices/internal/ui"
	"cli-view-crypto-prices/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	start, err := ui.RunIntro()
	if err != nil {
		log.Fatal(err)
	}

	if !start {
		os.Exit(0)
	}

	for {
		currency, limit, err := ui.RunForm()
		if err != nil {
			log.Fatal(err)
		}

		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			log.Fatal(err)
		}

		m := ui.NewModel(currency, limitInt)
		p := tea.NewProgram(m)

		finalModel, err := p.Run()
		if err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

		tableModel := finalModel.(ui.Model)
		if !tableModel.ShowForm {
			break
		}
		utils.ClearScreen()
	}
}
