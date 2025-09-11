package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	apiBaseURL = "https://api.coingecko.com/api/v3"
)

type Crypto struct {
	ID            string    `json:"id"`
	Symbol        string    `json:"symbol"`
	Name          string    `json:"name"`
	Image         string    `json:"image"`
	CurrentPrice  float64   `json:"current_price"`
	MarketCap     int64     `json:"market_cap"`
	MarketCapRank int       `json:"market_cap_rank"`
	LastUpdated   time.Time `json:"last_updated"`
	Description   string    `json:"description,omitempty"`
}

type CoinDescription struct {
	ID          string `json:"id"`
	Description *struct {
		En string `json:"en"`
	} `json:"description"`
}

func GetMarketData(currency string, perPage int) ([]Crypto, error) {
	url := fmt.Sprintf("%s/coins/markets?vs_currency=%s&order=market_cap_desc&per_page=%d&page=1&sparkline=false&precision=2", apiBaseURL, currency, perPage)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch market data: %s", resp.Status)
	}

	var cryptos []Crypto
	if err := json.NewDecoder(resp.Body).Decode(&cryptos); err != nil {
		return nil, err
	}

	return cryptos, nil
}

func GetCoinDescription(id string) (string, error) {
	url := fmt.Sprintf("%s/coins/%s", apiBaseURL, id)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusTooManyRequests {
			return "", fmt.Errorf("you have been rate-limited. Please wait a moment")
		}
		return "", fmt.Errorf("failed to fetch coin description: %s", resp.Status)
	}

	var coinDesc CoinDescription
	if err := json.NewDecoder(resp.Body).Decode(&coinDesc); err != nil {
		return "", err
	}

	if coinDesc.Description == nil || coinDesc.Description.En == "" {
		return "No description available.", nil
	}

	return coinDesc.Description.En, nil
}
