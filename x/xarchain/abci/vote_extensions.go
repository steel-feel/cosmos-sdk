package abci

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/log"
	"golang.org/x/sync/errgroup"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
	"xarchain/x/xarchain/keeper"
	
)

type VoteExtHandler struct {
	logger          log.Logger
	Height          int64               // current block height
	lastPriceSyncTS time.Time           // last time we synced prices
	providerTimeout time.Duration       // timeout for fetching prices from providers
	providers       map[string]Provider // mapping of chain-id to provider (e.g. 421614 -> Arbitrum Sepolia eth client)

	Keeper keeper.Keeper
}

/*
	CAVoteExtension struct

defines the canonical vote extension structure.
this is the object that will be marshaled as bytes and signed by the validator.
*/
type CAVoteExtension struct {
	Height    int64
	EventData map[string]IntentData // map of chainId to intent data
}

type IntentData struct {
	From    int64
	To      int64
	IDs     []string
}

func NewCAExtHandler(
	logger log.Logger, 
	keeper keeper.Keeper,
	Timeout time.Duration,
	supportedProviders map[string]string, /// NOTE: this is a map of chain-id to provider RPC URL

) *VoteExtHandler {
	return &VoteExtHandler{
		logger: logger,
		Keeper: keeper,
		providers: NewProvider(supportedProviders),
		providerTimeout: Timeout,
	}
}

func (h *VoteExtHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *abci.RequestExtendVote) (*abci.ResponseExtendVote, error){
			h.Height = req.Height
			h.lastPriceSyncTS = time.Now()
		
			g := new(errgroup.Group)
			providerAgg := NewProviderAggregator()
		
			var injectedVoteExtTx VoteExtensionTransaction
			// fetch previous from/to block numbers from txns
			if req.Height > 0 {
				for _, txn := range req.Txs {
					if err := json.Unmarshal(txn, &injectedVoteExtTx); err != nil {
						continue
					}
				}
			} 

			for chainId, provider := range h.providers {
				
				// Launch a goroutine to fetch events from provider.
				// Recall, vote extensions are not required to be deterministic.
				g.Go(func() error {
					doneCh := make(chan bool, 1)
					errCh := make(chan error, 1)
	
					var (
						intents map[string]EventsResp
						err    error
					)

					intents = make(map[string]EventsResp,3)
	
					go func() {
						
						tData, err1 := provider.FetchEvents(injectedVoteExtTx.IntentData[chainId].From, injectedVoteExtTx.IntentData[chainId].To)
						if err1 != nil {
							h.logger.Error("failed to fetch events from chain provider", "chainId", chainId, "err", err)
							errCh <- err
						}
						h.logger.Debug("Fetched events from chain provider", "chainId", chainId, "events", tData)

						intents[chainId] = tData

						doneCh <- true
					}()
	
					select {
					case <-doneCh:
						break
	
					case err := <-errCh:
						return err
	
					case <-time.After(h.providerTimeout):
						return fmt.Errorf("provider of chain %s timed out", chainId )
					}
	
					// aggregate and collect prices based on the base currency per provider
					for chainID, iData := range intents {
						success := providerAgg.SetIntentData(chainID, iData)
						if !success {
							return fmt.Errorf("failed to find any exchange rates in provider responses")
						}
					}
	
					return nil
				})
			}
	
			if err := g.Wait(); err != nil {
				// We failed to get some or all event from providers. In the case that
				// all events fail, we will throw error.
				h.logger.Error("failed to fetch events", "err", err)
			}
	
			// produce a canonical vote extension
			voteExt := CAVoteExtension{
				Height: req.Height,
				EventData: providerAgg.providerEvents,
			}
		
			// NOTE: We use stdlib JSON encoding, but an application may choose to use
			// a performant mechanism.
			bz, err := json.Marshal(voteExt)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal vote extension: %w", err)
			}
	
			return &abci.ResponseExtendVote{VoteExtension: bz}, nil
		}
}

func (h *VoteExtHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
	return func(ctx sdk.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
		var voteExt CAVoteExtension
		err := json.Unmarshal(req.VoteExtension, &voteExt)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal vote extension: %w", err)
		}

		if voteExt.Height != req.Height {
			return nil, fmt.Errorf("vote extension height does not match request height; expected: %d, got: %d", req.Height, voteExt.Height)
		}

		//TODO: Perform checks on events
		
		return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
	}
}
