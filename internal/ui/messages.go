package ui

import "cli-view-crypto-prices/internal/api"

// message for when the initial data is loaded
type dataLoadedMsg struct{ cryptos []api.Crypto }

// message for when description is fetched
type descriptionMsg struct{ description string }

// message for error during fetch
type errMsg struct{ err error }
