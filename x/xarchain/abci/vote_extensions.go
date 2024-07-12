package abci

import (
	"encoding/json"
	"fmt"
	"time"

	"cosmossdk.io/log"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"xarchain/x/xarchain/keeper"
	
)

type VoteExtHandler struct {
	logger          log.Logger
	currentBlock    int64                            // current block height
	lastPriceSyncTS time.Time                        // last time we synced prices
	Keeper keeper.Keeper
}

/* CAVoteExtension struct
defines the canonical vote extension structure.
this is the object that will be marshaled as bytes and signed by the validator.
*/
type CAVoteExtension struct {
	Height uint64
	Blocknumber uint64
	IDs []uint64
}

func (h *VoteExtHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
    return func(ctx sdk.Context, req *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {
		count := h.Keeper.GetTaskCount(ctx)
		var IDs []uint64
		var lastBlockNumeber uint64
		
		for i := 0; i < int(count); i++ { 
			task,found := h.Keeper.GetTask(ctx, uint64(i))
			if !found {
				continue
			}
			if task.Status != "picked"  {
				continue
			}

			currBlock, err := FetchTxn(task.Title)
			if err != nil {
				continue
			}

			if currBlock < 2 {
				continue
			}
			lastBlockNumeber = currBlock
			IDs = append(IDs, task.Id)
		} 
	
        voteExt := CAVoteExtension{
            IDs: IDs,
			Blocknumber: lastBlockNumeber,
			Height: uint64(req.Height),
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

		for id := range voteExt.IDs {
			task,found := h.Keeper.GetTask(ctx, uint64(id))
			if !found {
				return nil, fmt.Errorf("failed to find task id: %w", id)

			}
			if task.Status != "picked"  {
				return nil, fmt.Errorf("task is not picked yet: %w", id)

			}

			currBlock, err := FetchTxn(task.Title)
			if err != nil {
				return nil, fmt.Errorf("txn hash not found: %w", id)
			}

			if currBlock < 2 {
				return nil, fmt.Errorf("txn not mined yet: %w", id)
			}

		}

        return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
    }
}