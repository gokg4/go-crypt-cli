package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"cli-view-crypto-prices/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
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
		clearScreen()
	}
}
