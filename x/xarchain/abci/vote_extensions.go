package abci

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/log"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"os"
	"time"
	"xarchain/x/xarchain/keeper"
	"xarchain/x/xarchain/types"
)

type VoteExtHandler struct {
	logger log.Logger // current block height
	Keeper keeper.Keeper
}

/*
	CAVoteExtension struct

defines the canonical vote extension structure.
this is the object that will be marshaled as bytes and signed by the validator.
*/
type CAVoteExtension struct {
	Height      uint64
	Blocknumber int64
	IDs         []uint64
	TxHashs     []string
}

func NewCAExtHandler(
	logger log.Logger, // current block height             // last time we synced prices
	keeper keeper.Keeper,
) *VoteExtHandler {
	return &VoteExtHandler{
		logger: logger,
		Keeper: keeper,
	}
}

func (h *VoteExtHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {
		lastBlock, found := h.Keeper.GetCblock(ctx)
		if !found {
			return nil, fmt.Errorf("failed to get last block")
		}

		beforeEvent := time.Now()

		eventResp, err := FetchEvents(lastBlock.Blocknumber)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch events: %w", err)
		}

		afterEvent := time.Now()

		var IDs []uint64
		var TxHashs []string
		for _, intent := range eventResp.Intents {
			IDs = append(IDs, uint64(intent.ID))
			TxHashs = append(TxHashs, intent.TxHash)
		}

		f, err := os.OpenFile("/Users/himank/voteExt.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		// f.Write([]byte(fmt.Sprintf("Time taken to fetch events: %v, no of events %v \n", afterEvent.Sub(beforeEvent), len(eventResp.Intents) )) )
		f.WriteString(fmt.Sprintf("Blocknumber %v, Time taken to fetch events: %v, no of events %v \n", req.Height , afterEvent.Sub(beforeEvent), len(eventResp.Intents) ))
		f.Sync()
		defer f.Close()
		
		// defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), "xarchai2" )

		
		voteExt := CAVoteExtension{
			IDs:         IDs,
			Blocknumber: eventResp.lastBlock,
			Height:      uint64(req.Height),
			TxHashs:     TxHashs,
		}

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

		if voteExt.Height != uint64(req.Height) {
			return nil, fmt.Errorf("vote extension height does not match request height; expected: %d, got: %d", req.Height, voteExt.Height)
		}

		if len(voteExt.IDs) != len(voteExt.TxHashs) {
			return nil, fmt.Errorf("vote extension IDs and TxHashs length mismatch")
		}
		beforeEvent := time.Now()
		//code for loop to len(voteExt.IDs)

		for i := 0; i < len(voteExt.IDs); i++ {
			intent, found := h.Keeper.GetIntentById(ctx, uint64(voteExt.IDs[i]))
			if !found {
				return nil, fmt.Errorf("failed to find task id: %v", voteExt.IDs[i])
			}

			if intent.Status == "verified" {
				return nil, fmt.Errorf("intent is already verified: %v", voteExt.IDs[i])
			}

			txnDtls, err := FetchTxDetails(voteExt.TxHashs[i])
			if err != nil {
				return nil, fmt.Errorf("failed to fetch txn details: %w", err)
			}

			if txnDtls.To != "0x9ddB44C124E3e01D43ECEc91DD87B0BC9c4291FE" {
				return nil, fmt.Errorf("failed to To address is wrong: %w", err)
			}

			//NOTE: we can verify data also for additional checks
			// if txnDtls.Data != "0x" {
			// 	return nil, fmt.Errorf("failed to fetch data: %w", err)
			// }
		}

		afterEvent := time.Now()
		f, err := os.OpenFile("/Users/himank/verifyVote.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()
		f.WriteString(fmt.Sprintf("Blocknumber %v, Time taken to fetch events: %v, no of events %v \n", req.Height , afterEvent.Sub(beforeEvent), len(voteExt.IDs) ))
		f.Sync()
		
		defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), "xarchai2" )

		return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
	}
}
