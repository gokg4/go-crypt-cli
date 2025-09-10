package marketdata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	Description   string
}

func GetMarketData(c string, l string) (GeckoMarketData, error) {
	currency := strings.Trim(c, " ")
	topList := strings.Trim(l, " ")
	url := GeckoMarketUrl + fmt.Sprintf("?vs_currency=%v&order=market_cap_desc&per_page=%v&page=1&sparkline=false&precision=2", currency, topList)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		var data GeckoMarketData
		err := json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		return nil, fmt.Errorf("failed to fetch market data: status code %d", res.StatusCode)
	}
}

func FetchCoinDescription(id string) (string, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s", id)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Description map[string]string `json:"description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	desc, ok := result.Description["en"]
	if !ok {
		return "No description available.", nil
	}
	return desc, nil
}
