package abci

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func FetchTxn(txHashStr string) (uint64, error) {
	clientURL := "https://rpc.ankr.com/arbitrum_sepolia"

	client, err := ethclient.Dial(clientURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	txHash := common.HexToHash(txHashStr)

	ctx := context.Background()
	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		if err == ethereum.NotFound {
			return 0, nil
		} else {
			log.Fatalf("Failed to get transaction receipt: %v", err)
		}
		return 0, err
	}

	// Check if the block number is set, indicating the transaction has been mined
	if receipt.BlockNumber != nil {
		fmt.Printf("Transaction has been mined in block %d\n", receipt.BlockNumber.Uint64())
	} else {
		fmt.Println("Transaction has not been mined yet")
	}

	currNumber, err := client.BlockNumber(ctx)
	if err != nil {
		if err == ethereum.NotFound {
			return 1, nil
		} else {
			log.Fatalf("Failed to get transaction receipt: %v", err)
		}
		return 0, err
	}

	return currNumber, nil
}
