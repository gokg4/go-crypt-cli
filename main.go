package main

import (
	"gocryptocli/core"
	"gocryptocli/table"
)

func runApp() bool {
	data, currency := core.RestartApp()
	return table.CreateTable(data, currency)
}

func main() {
	for {
		shouldRestart := runApp()
		if !shouldRestart {
			break
		}
	}
}
