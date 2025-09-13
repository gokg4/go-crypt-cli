package ui

import (
	"cli-view-crypto-prices/internal/api"
	"fmt"
	"strings"
	"time"
)

func GenerateMarkdown(coin api.Crypto, currency string) (string, error) {
	var builder strings.Builder

	// H1 Header
	builder.WriteString(fmt.Sprintf("# Crypto Details %s\n\n", strings.ToUpper(currency)))

	// Details Table
	builder.WriteString("| Name | Symbol | Current Price | Market Cap | Market Cap Rank |\n")
	builder.WriteString("|---|---|---|---|---|\n")
	builder.WriteString(fmt.Sprintf("| %s | %s | %.2f %s | %d %s | %d |\n\n",
		coin.Name,
		strings.ToUpper(coin.Symbol),
		coin.CurrentPrice,
		strings.ToUpper(currency),
		coin.MarketCap,
		strings.ToUpper(currency),
		coin.MarketCapRank))

	// H2 Description
	builder.WriteString("## Description\n\n")
	builder.WriteString(fmt.Sprintf("%s\n\n", coin.Description))

	// Footer
	builder.WriteString("---\n\n")
	builder.WriteString(fmt.Sprintf("Generated on: %s\n\n", time.Now().Format("Jan 2, 2006 | 3:04PM MST")))
	builder.WriteString(fmt.Sprintf("[View on CoinGecko](https://www.coingecko.com/en/coins/%s)\n", coin.ID))

	return builder.String(), nil
}
