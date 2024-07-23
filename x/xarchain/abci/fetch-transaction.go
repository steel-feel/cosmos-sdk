package abci

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	// "strings"

	"github.com/ethereum/go-ethereum"
	// "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TxnDetails struct {
	To       string
	Data     string
	IntentID uint64
}

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

func FetchTxDetails(txHashStr string) (TxnDetails, error) {
	var txnDetails TxnDetails
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
			return txnDetails, nil
		} else {
			log.Fatalf("Failed to get transaction receipt: %v", err)
		}
		return txnDetails, err
	}

	if receipt.Status != 1 {
		log.Fatalf("Transaction not mined or wrong: %v", err)
		return txnDetails, errors.New("Transaction not mined or wrong")
	}

	tx, isPending, err := client.TransactionByHash(ctx, txHash)

	if err != nil {
		log.Fatalf("Failed to get transaction: %v", err)
		return txnDetails, err
	}

	if isPending {
		log.Fatalf("Transaction is pending")
		return txnDetails, errors.New("Transaction is pending")
	}

	txnDetails.To = tx.To().String()
	txnDetails.Data = fmt.Sprintf("%x", tx.Data())

	// we can decode logs and make verification on them as well like txHash refer to intent ID
	// receipt.Logs
	return txnDetails, nil
}

type EmittedIntents struct {
	ID     uint
	TxHash string
	Filer  string
}

type EventsResp struct {
	Intents   []EmittedIntents
	lastBlock int64
}

func GetToBlock(lastBlock int64, currBlockNumber int64) int64 {
	if lastBlock+999 > currBlockNumber {
		return currBlockNumber
	} else {
		return lastBlock + 999
	}
}

func FetchEvents(lastBlock int64) (*EventsResp, error) {
	clientURL := "https://rpc.ankr.com/arbitrum_sepolia"

	client, err := ethclient.Dial(clientURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
		return nil, err

	}

	// Get the current block number
	currBlockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("Failed to get the block number: %v", err)
		return nil, err

	}

	ToBlock := GetToBlock(lastBlock, int64(currBlockNumber))

	//to block should be the current block number or +999 whichever lower
	contractAddress := common.HexToAddress("0x9ddB44C124E3e01D43ECEc91DD87B0BC9c4291FE") // Replace with your contract address
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(lastBlock),
		ToBlock:   big.NewInt(ToBlock),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to retrieve logs: %v", err)
		return nil, err

	}

	var emittedIntents []EmittedIntents

	for _, vLog := range logs {
		event := struct {
			Filer    common.Address 
			IntentID *big.Int   
		}{}

		event.Filer = common.HexToAddress(vLog.Topics[1].Hex())
        event.IntentID = new(big.Int)
        event.IntentID.SetString(vLog.Topics[2].Hex()[2:], 16)

		fmt.Printf("Log: %+v\n", event)
		emittedIntents = append(emittedIntents, EmittedIntents{
			TxHash: vLog.TxHash.String(),
			Filer:  event.Filer.String(),
			ID:     uint(event.IntentID.Uint64()),
		})
	}

	return &EventsResp{
		Intents:   emittedIntents,
		lastBlock: ToBlock,
	}, nil

}
