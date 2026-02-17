# Scoville Trading Bot 🌶️

Scoville is a terminal-based application for running and monitoring a phased, simulated trading strategy. It is designed to model a market-making or price-support scenario through a sequence of automated steps, providing real-time feedback through an interactive TUI dashboard.

## Features

-   **Phased Trading Strategy:** Simulates market activity in distinct, sequential phases to systematically achieve a target price.
-   **Interactive TUI:** A dashboard built with `bubbletea` provides real-time feedback on the bot's progress, including the current phase, completion percentage, and a detailed log of actions.
-   **State Persistence:** The bot's progress is saved to `scoville_progress.json`, allowing the simulation to be stopped and resumed from where it left off.
-   **Live & Paper Modes:** Run in a safe simulation mode (`paper mode`) without spending real funds, or connect to a live network to execute real trades.
-   **Configurable:** All strategy parameters, network settings, and wallet credentials are controlled via a simple `.env` file.

## How It Works

The bot executes a strategy divided into three main phases:

1.  **🐋 Phase 1: Whale Anchoring**
    This phase simulates a small number of large buys from a single wallet. The goal is to create significant price movement and establish a psychological price floor, appearing as if a large investor ("whale") is taking a position.

2.  **🐟 Phase 2: Retail Velocity**
    Following the whale buys, this phase simulates a larger number of smaller, more frequent buys. This activity is designed to mimic retail investor interest, filling in gaps in the price chart and creating a sense of organic momentum.

3.  **🔒 Phase 3: Liquidity Provision (Optional)**
    This phase is controlled by the `ENABLE_LIQUIDITY_PHASE` setting. If enabled, it focuses on establishing and supporting the token's liquidity on a decentralized exchange to ensure market stability at the new price level.

## Getting Started

### Prerequisites

-   [Go](https://golang.org/doc/install) (version 1.21 or newer)

### Installation & Setup

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/chili-network/scoville.git
    cd scoville
    ```

2.  **Configure your environment:**
    Copy the example environment file and edit it with your own settings.
    ```sh
    cp .env.example .env
    ```
    See the [Configuration Details](#configuration-details) section below for a full explanation of each variable.

3.  **Install dependencies:**
    The `Makefile` simplifies the process. This command will download and verify the necessary Go modules.
    ```sh
    make install
    ```

### Running the Bot

-   **To run in Paper Mode (Simulation):**
    This is the default and recommended mode for testing. It will simulate all trades without connecting to a wallet or spending real funds.
    ```sh
    make run
    ```

-   **To run in Live Mode (Real Trades):**
    This will connect to the specified RPC, use the provided private key, and execute real trades on the blockchain.
    ```sh
    make run-live
    ```
    You will be prompted with a warning and must type `YES` to confirm before the bot starts.

## TUI Controls

While the application is running, you can use the following keys:

-   **`P`**: Pause or Resume the current operation. The bot will halt its activity and wait for you to resume.
-   **`Q`**: Initiate the quit sequence. You will be asked to confirm.
-   **`Y` / `N`**: Confirm (`Y`) or cancel (`N`) quitting the application.

---

## Configuration Details

All configuration is handled in the `.env` file.

| Variable                  | Description                                                                                                                              | Example                                                   |
| ------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------- |
| **RPC_URL**               | Your connection to the blockchain. Get a free key from a provider like Infura, Alchemy, or QuickNode.                                      | `https://mainnet.infura.io/v3/YOUR_API_KEY_HERE`          |
| **CHAIN_ID**              | The ID of the target blockchain (e.g., 1 for Ethereum Mainnet, 8453 for Base, 137 for Polygon).                                            | `1`                                                       |
| **PRIVATE_KEY**           | **CRITICAL:** The private key of the wallet that will execute the trades. For security, this **must** be a fresh wallet funded only with the necessary ETH (for gas) and stablecoins (for buys). **DO NOT** use a personal or deployer wallet. | `0000...0000`                                             |
| **TOKEN_ADDRESS**         | The contract address of the token you are targeting.                                                                                       | `0x83e8fb8d8176224fcc828edc73e152ec1818a2da`              |
| **ROUTER_ADDRESS**        | The contract address of the DEX router (e.g., Uniswap V2, SushiSwap).                                                                    | `0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D`              |
| **TARGET_PRICE**          | The target price in USD that you want the token to reach. The simulation will run until this goal is met or the budget is exhausted.      | `0.00075`                                                 |
| **TOTAL_BUDGET**          | The total amount of USD (in stablecoins) allocated to this entire strategy.                                                                | `3000.0`                                                  |
| **WHALE_BUY_MIN**         | The minimum USD value for a single Phase 1 "whale" buy.                                                                                    | `200.0`                                                   |
| **WHALE_BUY_MAX**         | The maximum USD value for a single Phase 1 "whale" buy. The actual amount will be a random value between min and max.                       | `500.0`                                                   |
| **RETAIL_BUY_MIN**        | The minimum USD value for a single Phase 2 "retail" buy.                                                                                 | `20.0`                                                    |
| **RETAIL_BUY_MAX**        | The maximum USD value for a single Phase 2 "retail" buy. The actual amount will be a random value between min and max.                    | `50.0`                                                    |
| **ENABLE_LIQUIDITY_PHASE**| If `false`, the bot will increase Phase 2 buys and finish. If `true`, the bot will proceed to Phase 3.                                   | `false`                                                   |
| **PAPER_MODE**            | Set to `true` to simulate all trades (no real funds spent). Set to `false` to enable live trading. The `make run` commands manage this for you. | `true`                                                    |

---

## ⚠️ Disclaimer

**This software is for educational and experimental purposes only. Live trading involves significant financial risk, including the potential loss of all invested funds. The authors and contributors are not responsible for any financial losses you may incur. Always do your own research and never risk more than you are willing to lose.**
