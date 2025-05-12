# Indexer

Usage example:

```go
package main

import (
    "context"
    "time"

    "github.com/ethereum/go-ethereum/core/types"
    "github.com/xmtp/xmtpd/pkg/indexerpoc"
    "github.com/xmtp/xmtpd/pkg/indexerpoc/contract"
    "github.com/xmtp/xmtpd/pkg/indexerpoc/contract/group_message_broadcaster"
    "github.com/xmtp/xmtpd/pkg/indexerpoc/contract/identity_update_broadcaster"
    "go.uber.org/zap"
)

func main() {

ctx := context.Background()

storage := indexerpoc.NewMemoryStorage()
log, _ := zap.NewProduction()
defer log.Sync()

// Define networks.
ethereumNetwork := indexerpoc.NewNetwork(
    "ethereum",
    1,
    "https://mainnet.infura.io/v3/YOUR_API_KEY",
    12 * time.Second,
)

arbitrumNetwork := indexerpoc.NewNetwork(
    "arbitrum",
    42161,
    "https://arb1.arbitrum.io/rpc",
    250 * time.Millisecond,
)

// Create indexer.
manager, err := indexerpoc.NewIndexer(
        ctx,
        log,
        storage,
        5*time.Second,
        10,
        4,
        []*indexerpoc.Network{ethereumNetwork, arbitrumNetwork}
    )
if err != nil {
    log.Error("Failed to create indexer", zap.Error(err))
    return
}

// Generic log processor.
// This function would hold same logic as StoreLog for group message and identity updates.
createLogProcessor := func(networkName string) contract.LogProcessor {
    loggerCopy := log
    return func(ctx context.Context, logs []types.Log) error {
        for _, event := range logs {
            loggerCopy.Info(networkName+" Transfer event",
                zap.String("from", event.Topics[1].Hex()),
                zap.String("to", event.Topics[2].Hex()),
                zap.Binary("data", event.Data),
            )
        }
        return nil
    }
}

// Generic reorg processor.
simpleDeleteReorgProcessor := func(ctx context.Context, fromBlock uint64) error {
    log.Warn("Simple delete reorg handler",
        zap.Uint64("fromBlock", fromBlock),
    )

    return nil
}

// Create processors for each network.
ethereumTransferProcessor := createLogProcessor("Ethereum")
arbitrumTransferProcessor := createLogProcessor("Arbitrum")

// ERC20 Transfer event signature has the same signature across networks.
transferEventSignature := "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

// Add contracts to index with different reorg handling approaches, starting blocks, etc.
contracts := []contract.Contract{
    contract.NewContract(
        "USDC_Ethereum",
        1, // Ethereum mainnet
        "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
        15000000,
        []string{transferEventSignature},
        ethereumTransferProcessor,
        simpleDeleteReorgProcessor,
    ),
    contract.NewContract(
        "USDC_Arbitrum",
        42161, // Arbitrum One
        "0xff970a61a04b1ca14834a43f5de4533ebddb5cc8",
        7000000,
        []string{transferEventSignature},
        arbitrumTransferProcessor,
        simpleDeleteReorgProcessor,
    ),
}

// Add all contracts to be indexed
if err := manager.AddContracts(contracts); err != nil {
    log.Error("Failed to add contracts", zap.Error(err))
    return
}

// Start indexing.
manager.Run()
}
```

## Using Specialized Contract Implementations

For more complex use cases, you can create specialized contract implementations:

```go
package main

import (
    "context"
    "database/sql"
    "time"

    "github.com/xmtp/xmtpd/pkg/db/queries"
    "github.com/xmtp/xmtpd/pkg/indexerpoc"
    "github.com/xmtp/xmtpd/pkg/indexerpoc/contract/group_message_broadcaster"
    "go.uber.org/zap"
)

func main() {
    ctx := context.Background()
    storage := indexerpoc.NewMemoryStorage()
    log, _ := zap.NewProduction()
    defer log.Sync()

    // Define networks
    mainnet := indexerpoc.NewNetwork(
        "ethereum",
        1,
        "https://mainnet.infura.io/v3/YOUR_API_KEY",
        12 * time.Second,
    )

    // Create indexer
    manager, err := indexerpoc.NewIndexer(
        ctx,
        log,
        storage,
        5*time.Second,
        10,
        4,
        []*indexerpoc.Network{mainnet},
    )
    if err != nil {
        log.Error("Failed to create indexer", zap.Error(err))
        return
    }

    // Setup database connection
    db, err := sql.Open("postgres", "your_db_connection_string_here")
    if err != nil {
        log.Error("Failed to connect to database", zap.Error(err))
        return
    }
    defer db.Close()

    // Create querier
    querier := queries.New(db)

    // Create specialized contracts
    groupMessageContract := group_message_broadcaster.NewGroupMessageContract(
        log,           // Logger
        querier,       // Database querier
        "GroupMessage", // Contract name
        1,             // Chain ID (Ethereum mainnet)
        "0x1234567890abcdef1234567890abcdef12345678", // Contract address
        15000000,      // Start block
    )

    // GroupMessageContract now implements the indexerpoc.Contract interface directly
    // Add contracts to indexer
    if err := manager.AddContracts([]indexerpoc.Contract{
        groupMessageContract, // GroupMessageContract itself implements the interface
    }); err != nil {
        log.Error("Failed to add contracts", zap.Error(err))
        return
    }

    // Start indexing
    manager.Run()
}
```
