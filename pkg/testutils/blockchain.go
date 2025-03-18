package testutils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func Int64ToHash(x int64) common.Hash {
	return common.BigToHash(big.NewInt(x))
}

func LogTransactionsIfFailed(t *testing.T, rpcURL string) {
	if t.Failed() {
		t.Log("Failed test. Logging transactions...")
		LogAllTransactions(t, rpcURL)
		InspectMempool(t, rpcURL)
	}
}

func LogAllTransactions(t *testing.T, rpcURL string) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		t.Fatalf("Failed to connect to Anvil: %v", err)
	}
	defer client.Close()

	var logOutput string

	ctx := context.Background()

	// Get the latest block number
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get latest block header: %v", err)
	}
	latestBlock := header.Number

	logOutput += fmt.Sprintf("🔹 Fetching transactions from block 0 to %d...", latestBlock)

	// Iterate through all blocks from genesis to latest
	for i := int64(0); i <= latestBlock.Int64(); i++ {
		block, err := client.BlockByNumber(ctx, big.NewInt(i))
		if err != nil {
			logOutput += fmt.Sprintf("Failed to get block %d: %v", i, err)
			continue
		}

		logOutput += fmt.Sprintf(
			"\n📌 Block %d (Hash: %s) - Transactions: %d",
			i,
			block.Hash().Hex(),
			len(block.Transactions()),
		)

		// Iterate through all transactions in the block
		for _, tx := range block.Transactions() {
			receipt, err := client.TransactionReceipt(ctx, tx.Hash())
			if err != nil {
				logOutput += fmt.Sprintf("\n  Tx %s: Pending", tx.Hash().Hex())
			} else {
				status := "Success"
				if receipt.Status == 0 {
					status = "Failed"
				}
				logOutput += fmt.Sprintf("\n  Tx %s: Mined in Block %d (Status: %s, Nonce: %d)",
					tx.Hash().Hex(), receipt.BlockNumber.Uint64(), status, tx.Nonce())
			}
		}
	}

	// Fetch and log all pending transactions
	logOutput += fmt.Sprintf("\n🔹 Fetching pending transactions...")
	pendingBlock, err := client.BlockByNumber(ctx, big.NewInt(-1)) // -1 for "pending"
	if err != nil {
		logOutput += fmt.Sprintf("\nFailed to get pending block: %v", err)
	} else if len(pendingBlock.Transactions()) == 0 {
		logOutput += fmt.Sprintf("\nNo transactions pending.")
	} else {
		for _, tx := range pendingBlock.Transactions() {
			logOutput += fmt.Sprintf("\n  Pending Tx: %s (o: %s, Gas: %d, Value: %s, Nonce: %d)",
				tx.Hash().Hex(),
				func() string {
					if tx.To() != nil {
						return tx.To().Hex()
					}
					return "Contract Creation"
				}(),
				tx.Gas(),
				tx.Value().String(),
				tx.Nonce(),
			)
		}
	}

	t.Log(logOutput)
}

func SubscribePendingTxs(t *testing.T, rpcURL string) error {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	// Connect to WebSocket endpoint
	client, err := rpc.DialContext(context.Background(), strings.Replace(rpcURL, "http", "ws", 1))
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum WS: %v", err)
	}
	defer client.Close()

	// Channel to receive transaction hashes
	txHashes := make(chan string)

	// Subscribe to pending transactions
	sub, err := client.Subscribe(context.Background(), "eth", txHashes, "newPendingTransactions")
	if err != nil {
		return fmt.Errorf("failed to subscribe to pending transactions: %v", err)
	}

	// Listen for incoming transactions
	fmt.Println("Listening for new pending transactions...")
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-sub.Err():
			t.Logf("Subscription error: %v\n", err)
			return err
		case txHash := <-txHashes:
			t.Logf("New Pending Transaction: %s\n", txHash)
		}
	}
}

func SubscribeToLogs(t *testing.T, rpcURL string, contractAddress string) error {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	client, err := ethclient.DialContext(
		context.Background(),
		strings.Replace(rpcURL, "http", "ws", 1),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum WS: %v", err)
	}
	defer client.Close()

	// Channel to receive logs
	logs := make(chan types.Log)

	// Define log filter for a specific contract
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(contractAddress)},
	}

	// Subscribe to contract events
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %v", err)
	}

	fmt.Printf("Listening for logs from contract: %s\n", contractAddress)
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-sub.Err():
			t.Logf("Subscription error: %v\n", err)
			return err
		case logEntry := <-logs:
			t.Logf("New Log from %s: %+v\n", logEntry.Address.Hex(), logEntry)
		}
	}
}

// TxpoolStatus represents the response from `txpool_status`
type TxpoolStatus struct {
	Pending string `json:"pending"`
	Queued  string `json:"queued"`
}

// TxpoolInspect represents the response from `txpool_inspect`
type TxpoolInspect map[string]map[string]map[string]string

// TxpoolContent represents the response from `txpool_content`
type TxpoolContent struct {
	Pending map[string]map[string]TransactionInfo `json:"pending"`
	Queued  map[string]map[string]TransactionInfo `json:"queued"`
}

// TransactionInfo holds details about a transaction in the txpool
type TransactionInfo struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Value    string `json:"value"`
	Nonce    string `json:"nonce"`
	Input    string `json:"input"`
}

// InspectMempool fetches mempool details from Anvil using txpool_status, txpool_inspect, and txpool_content
func InspectMempool(t *testing.T, rpcURL string) {
	// Connect to the Ethereum JSON-RPC endpoint
	client, err := rpc.DialContext(context.Background(), rpcURL)
	require.NoError(t, err)
	defer client.Close()

	var logOutput string

	// Fetch txpool_status
	var status TxpoolStatus
	err = client.Call(&status, "txpool_status")
	require.NoError(t, err)
	logOutput += fmt.Sprintf(
		"\n🔹 Mempool Status: Pending: %s | Queued: %s\n",
		status.Pending,
		status.Queued,
	)

	// Fetch txpool_inspect
	var inspect TxpoolInspect
	err = client.Call(&inspect, "txpool_inspect")
	require.NoError(t, err)
	logOutput += fmt.Sprintf("\n🔹 Mempool Inspection:\n")
	for from, txs := range inspect {
		logOutput += fmt.Sprintf("From: %s\n", from)
		for account, details := range txs {
			logOutput += fmt.Sprintf("  Account: %s -> \n", account)
			for nonce, tx := range details {
				txJSON, _ := json.MarshalIndent(tx, "  ", "  ")
				logOutput += fmt.Sprintf("    Nonce: %s\n    %s\n", nonce, string(txJSON))
			}
		}
	}

	// Fetch txpool_content
	var content TxpoolContent
	err = client.Call(&content, "txpool_content")
	require.NoError(t, err)

	logOutput += fmt.Sprintf("\n🔹 Detailed Mempool Transactions:")
	logOutput += fmt.Sprintf(formatTxpoolContent(t, "Pending Transactions", content.Pending))
	logOutput += fmt.Sprintf(formatTxpoolContent(t, "Queued Transactions", content.Queued))

	t.Log(logOutput)
}

// printTxpoolContent prints transactions in a structured way
func formatTxpoolContent(
	t *testing.T,
	title string,
	txMap map[string]map[string]TransactionInfo,
) string {
	var logOutput string
	logOutput += fmt.Sprintf("\n📌 %s:\n", title)

	for from, txs := range txMap {
		logOutput += fmt.Sprintf("From: %s\n", from)
		for nonce, tx := range txs {
			txJSON, _ := json.MarshalIndent(tx, "  ", "  ")
			logOutput += fmt.Sprintf("  Nonce: %s\n  %s\n", nonce, string(txJSON))
		}
	}
	return logOutput
}
