package abci

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"sync"
	// "strings"

	"github.com/ethereum/go-ethereum"
	// "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/crypto"
)

type TxnDetails struct {
	To       string
	Data     string
	IntentID uint64
}

type Provider struct {
	client          ethclient.Client
	contractAddress string
}

func NewProvider(supportedProviders map[string]string) map[string]Provider {
	providers := make(map[string]Provider)
	for chainID, rpcURL := range supportedProviders {
		// Supporting only Arbitrum Sepolia, OP Sepolia and Ethereum Sepolia
		if chainID == "421614" || chainID == "11155111" || chainID == "11155420" {
			client, err := ethclient.Dial(rpcURL)
			if err != nil {
				log.Fatalf("Failed to connect to the Ethereum client: %v", err)
			}

			switch chainID {
			case "421614":
				providers[chainID] = Provider{
					client:          *client,
					contractAddress: "0xF5620427CB929BAdd689f92D1AE52704dD019BDA",
				}

			case "11155111":
			

			case "11155420":
				providers[chainID] = Provider{
					client:          *client,
					contractAddress: "0x2884bD2cf67b933CBb5199093Cea052d7A79198A",
				}
			}
		}

	}
	return providers

}

func (p *Provider) FetchEvents(prevFrom int64, prevTo int64) (EventsResp, error) {
	var evtResponse EventsResp
	// Get the current block number
	currentBlock, err := p.client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("Failed to get the block number: %v", err)
		return evtResponse, err

	}

	iCurrentBlock := int64(currentBlock)

	if prevTo == iCurrentBlock {
		return evtResponse, nil
	}

	var From int64
	if prevFrom == 0 {
		From = iCurrentBlock
	} else {
		From = prevTo + 1
	}

	To := iCurrentBlock
	if From+999 < iCurrentBlock {
		To = From + 999
	}

    // Define the event signature (you can get this from the ABI or Etherscan)
	eventSignature := []byte("IntentFulfiled(address,bytes32)")
    eventSignatureHash := common.BytesToHash(crypto.Keccak256Hash(eventSignature).Bytes())


	contractAddress := common.HexToAddress(p.contractAddress)
	query := ethereum.FilterQuery{
        Topics:    [][]common.Hash{{eventSignatureHash}},
		FromBlock: big.NewInt(From),
		ToBlock:   big.NewInt(To),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := p.client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to retrieve logs: %v, From: %v, To: %v  ", err, From, To)
		return evtResponse, err
	}

	var emittedIntents []EmittedIntents
	for _, vLog := range logs {
		//TODO: Emitted event struct
		event := struct {
			Filer    common.Address
			IntentID [32]byte
		}{}

        event.Filer = common.BytesToAddress(vLog.Topics[1].Bytes())
        copy(event.IntentID[:], vLog.Topics[2].Bytes())

		event.Filer = common.HexToAddress(vLog.Topics[1].Hex())

		emittedIntents = append(emittedIntents, EmittedIntents{
			TxHash: vLog.TxHash.String(),
			Filer:  event.Filer.String(),
			ID:     string(event.IntentID[:]),
		})
	}

	return EventsResp{
		Intents: emittedIntents,
		To:      To,
		From:    From,
	}, nil

}

// // ProviderAggregator is a simple aggregator for provider prices and candles.
// // It is thread-safe since it is assumed to be called concurrently in price
// // fetching goroutines, i.e. ExtendVote.
type ProviderAggregator struct {
	mtx sync.Mutex

	providerEvents map[string]IntentData
}

func NewProviderAggregator() *ProviderAggregator {
	return &ProviderAggregator{
		providerEvents: make(map[string]IntentData),
	}
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
	ID     string
	TxHash string
	Filer  string
}

type EventsResp struct {
	Intents []EmittedIntents
	From    int64
	To      int64
}

func (p *ProviderAggregator) SetIntentData(chainId string, iEvent EventsResp) bool {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	var IDs []string
	
	for _, intent := range iEvent.Intents {
		IDs = append(IDs, intent.ID)
	}

	p.providerEvents[chainId] = IntentData{
		From:    iEvent.From,
		To:      iEvent.To,
		IDs:     IDs,
	}

	return true
}
