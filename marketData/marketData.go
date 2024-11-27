package marketdata

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const GeckoMarketUrl = "https://api.coingecko.com/api/v3/coins/markets"

type GeckoMarketData []struct {
	ID            string    `json:"id"`
	Symbol        string    `json:"symbol"`
	Name          string    `json:"name"`
	Image         string    `json:"image"`
	CurrentPrice  float64   `json:"current_price"`
	MarketCap     int64     `json:"market_cap"`
	MarketCapRank int       `json:"market_cap_rank"`
	LastUpdated   time.Time `json:"last_updated"`
}

func MarketData(c string, l string) GeckoMarketData {
	currency := strings.Trim(c, " ")
	topList := strings.Trim(l, " ")
	url := GeckoMarketUrl + fmt.Sprintf("?vs_currency=%v&order=market_cap_desc&per_page=%v&page=1&sparkline=false&precision=2", currency, topList)

	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode == 200 {
		var data GeckoMarketData // Assuming an array of market data
		err := json.Unmarshal(body, &data)
		if err != nil {
			log.Fatal(err)
		}

		return data

	} else {
		fmt.Printf("Error Code: %d \n", res.StatusCode)
		os.Exit(1)
		return nil
	}
}
