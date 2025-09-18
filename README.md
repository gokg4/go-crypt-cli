![GitHub](https://img.shields.io/github/license/gokg4/go-crypt-cli) ![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/gokg4/go-crypt-cli) ![GitHub repo size](https://img.shields.io/github/repo-size/gokg4/go-crypt-cli)

# Crypto Price Viewer CLI

A command-line interface (CLI) application for viewing cryptocurrency prices, built with Go and the Bubble Tea framework.

## Features

*   **Top Cryptocurrencies:** View a list of the top cryptocurrencies by market capitalization.
*   **Detailed View:** Select a cryptocurrency to view its detailed information, including a description.
*   **Customizable:** Choose the currency (e.g., USD, EUR, JPY) and the number of cryptocurrencies to display.
*   **Markdown Export:** Save the details of a cryptocurrency to a markdown file.
*   **Interactive UI:** A user-friendly interface built with Bubble Tea.

## How It Works

The application fetches cryptocurrency market data from the [CoinGecko API](https://www.coingecko.com/en/api). The data is then displayed in a table format. When a user selects a cryptocurrency, the application fetches the coin's description and displays it in a detailed view. The user can then save this information to a markdown file.

## Getting Started

### Prerequisites

*   Go 1.18 or later
*   A working internet connection

### Installation

1.  Clone the repository:
    ```bash
    git clone https://github.com/gokg4/go-crypt-cli
    ```
2.  Navigate to the project directory:
    ```bash
    cd go-crypt-cli
    ```
3.  Install the dependencies:
    ```bash
    go mod tidy
    ```
4.  Build and Compile with the following command.
    ```bash
    go build -o ./bin/geckoCrypto -ldflags="-s -w"
    ```

### Running the Application

To run the application, use the following command:

```bash
go run .
```

By default, the application will display the top 10 cryptocurrencies in USD. You can customize this by providing command-line flags:

*   `-currency`: The currency to display the prices in (e.g., 'eur', 'jpy'). Defaults to 'usd'.
*   `-limit`: The number of cryptocurrencies to display. Defaults to 10.

For example, to display the top 20 cryptocurrencies in EUR, you would run:

```bash
go run . -currency eur -limit 20
```

## Usage

Once built and compiled you can run from command line using the following command.

```bash
./bin/geckoCrypto
```

### Main View

*   **Up/Down Arrows:** Navigate through the list of cryptocurrencies.
*   **Enter:** View the details of the selected cryptocurrency.
*   **e:** Edit preferences (currency and limit). This will exit the application and you will have to re-run it with the desired flags.
*   **q / ctrl+c:** Quit the application.

### Details View

*   **Up/Down Arrows:** Scroll through the cryptocurrency description.
*   **m:** Save the details of the cryptocurrency to a markdown file. The file will be saved in the `markdown` directory.
*   **Enter / Esc:** Return to the main list view.

## Architecture

This project is built using the **Elm Architecture**, a design pattern for building interactive applications. The Go library it uses, `bubbletea`, is directly inspired by this pattern.

The core of the architecture is a unidirectional data flow:

```
+-----------------+      +----------------+      +-----------------+
|      Model      |----->|       View     |----->|       UI        |
| (holds state)   |      | (renders HTML) |      |   (the screen)  |
+-----------------+      +----------------+      +-----------------+
        ^                                                 |
        |                                                 | (User input, e.g. key press)
        |                                                 |
+-----------------+      +----------------+      +-----------------+
|      New        |<- - -|     Command    |      |     Message     |
|   Model & Cmd   |      |   (optional)   |      |  (e.g. KeyMsg)  |
+-----------------+      +----------------+      +-----------------+
        ^                                                 |
        |                                                 |
        +-------------------------------------------------+
        |                      Update                     |
        |           (processes messages, updates state)   |
        +-------------------------------------------------+
```

1.  **Model:** A single struct (`internal/ui/model.go`) holds the entire state of the application.
2.  **View:** A function (`internal/ui/view.go`) takes the model and returns a `string` to be rendered to the terminal. It is a pure representation of the state.
3.  **Update:** A function (`internal/ui/update.go`) is the only place where the state can be changed. It takes a message (like a key press or an API response) and the current model, and returns a new, updated model.
4.  **Commands:** When the `Update` function needs to perform an action that has side effects (like an HTTP request), it returns a `tea.Cmd`. The `bubbletea` runtime executes this command, which then sends a new message back to the `Update` function with the result.

This creates a clear and predictable data flow that makes the application easier to reason about, debug, and maintain.


## Dependencies

This project uses the following Go packages:

*   [github.com/charmbracelet/bubbles](https://github.com/charmbracelet/bubbles)
*   [github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)
*   [github.com/charmbracelet/glamour](https://github.com/charmbracelet/glamour)
*   [github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)
*   [github.com/muesli/reflow](https://github.com/muesli/reflow)
*   [github.com/spf13/viper](https://github.com/spf13/viper)

All dependencies are managed using Go modules.

## Contributors

- [gokg4](https://github.com/gokg4) - creator and maintainer

***

Made with [Charm](https://charm.sh/libs/).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400"></a>