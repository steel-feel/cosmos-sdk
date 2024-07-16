package abci

import (
	"encoding/json"
	"errors"
	"fmt"

	"cosmossdk.io/log"

	abci "github.com/cometbft/cometbft/abci/types"
	// cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"xarchain/x/xarchain/keeper"
)

// StakeWeightedPrices defines the structure a proposer should use to calculate
// and submit the stake-weighted prices for a given set of supported currency
// pairs, in addition to the vote extensions used to calculate them. This is so
// validators can verify the proposer's calculations.

// Its kind of temp storage in block formation
// before using this struct to commit the values to storage
type SuccessTransactionsID struct {
	TaskIDs            []uint64
	ExtendedCommitInfo abci.ExtendedCommitInfo
}

type ProposalHandler struct {
	logger   log.Logger
	keeper   keeper.Keeper          // our oracle module keeper
	valStore baseapp.ValidatorStore // to get the current validators' pubkeys
}

func NewProposalHandler(logger log.Logger, keeper keeper.Keeper, valStore baseapp.ValidatorStore) *ProposalHandler {
	return &ProposalHandler{
		logger:   logger,
		keeper:   keeper,
		valStore: valStore,
	}
}

func (h *ProposalHandler) PrepareProposal() sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
		proposalTxs := req.Txs

		if req.Height >= ctx.ConsensusParams().Abci.VoteExtensionsEnableHeight && ctx.ConsensusParams().Abci.VoteExtensionsEnableHeight != 0 {
			err := baseapp.ValidateVoteExtensions(ctx, h.valStore, req.Height, ctx.ChainID(), req.LocalLastCommit)
			if err != nil {
				return nil, err
			}

			goodTxns, err := h.computeCAIds(ctx, req.LocalLastCommit)
			if err != nil {
				return &abci.ResponsePrepareProposal{
					Txs: proposalTxs,
				}, nil
			}

			if len(goodTxns) == 0 {
				return nil, errors.New("no good transactions")
			}

			injectedVoteExtTx := SuccessTransactionsID{
				TaskIDs:            goodTxns,
				ExtendedCommitInfo: req.LocalLastCommit,
			}

			// NOTE: We use stdlib JSON encoding, but an application may choose to use
			// a performant mechanism. This is for demo purposes only.
			bz, err := json.Marshal(injectedVoteExtTx)
			if err != nil {
				h.logger.Error("failed to encode injected vote extension tx", "err", err)
				return nil, errors.New("failed to encode injected vote extension tx")
			}

			// Inject a "fake" tx into the proposal s.t. validators can decode, verify,
			// and store the canonical stake-weighted average prices.
			proposalTxs = append(proposalTxs, bz)
		}

		// proceed with normal block proposal construction, e.g. POB, normal txs, etc...

		return &abci.ResponsePrepareProposal{
			Txs: proposalTxs,
		}, nil
	}
}

func (h *ProposalHandler) ProcessProposal() sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		if len(req.Txs) == 0 {
			return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
		}
		if req.Height > ctx.ConsensusParams().Abci.VoteExtensionsEnableHeight {
			var injectedVoteExtTx SuccessTransactionsID
			if err := json.Unmarshal(req.Txs[0], &injectedVoteExtTx); err != nil {
				h.logger.Error("failed to decode proccess Proposal", "err", err)
				// return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
				return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
			}

			err := baseapp.ValidateVoteExtensions(ctx, h.valStore, req.Height, ctx.ChainID(), injectedVoteExtTx.ExtendedCommitInfo)
			if err != nil {
				return nil, err
			}
		}

		// // Verify the proposer's stake-weighted oracle prices by computing the same
		// // calculation and comparing the results. We omit verification for brevity
		// // and demo purposes.
		// goodTxns, err := h.computeCAIds(ctx, injectedVoteExtTx.ExtendedCommitInfo)
		// if err != nil {
		// 	return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
		// }
		// if err := compareOraclePrices(injectedVoteExtTx.StakeWeightedPrices, stakeWeightedPrices); err != nil {
		// 	return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
		// }

		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}
}

func (h *ProposalHandler) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	res := &sdk.ResponsePreBlock{}
	if len(req.Txs) == 0 {
		return res, nil
	}
	if req.Height > ctx.ConsensusParams().Abci.VoteExtensionsEnableHeight {
		var injectedVoteExtTx SuccessTransactionsID
		if err := json.Unmarshal(req.Txs[0], &injectedVoteExtTx); err != nil {
			h.logger.Error("failed to decode preblocker vote", "err", err)
			// return nil, err
			return res, nil

		}

		for _, taskID := range injectedVoteExtTx.TaskIDs {
			task, found := h.keeper.GetTask(ctx, taskID)
			if !found {
				h.logger.Error("failed to get Task %v", taskID)
				return nil, fmt.Errorf("failed to get Task %v", taskID)
			}

			task.Status = "verified"
			h.keeper.SetTask(ctx, task)

		}
	}

	return res, nil
}

func (h *ProposalHandler) computeCAIds(ctx sdk.Context, ci abci.ExtendedCommitInfo) ([]uint64, error) {
	var voteExt CAVoteExtension
	var IDs []uint64
	if len(ci.Votes) == 0 {
		return IDs, errors.New("no votes in commit info")
	}

	if err := json.Unmarshal(ci.Votes[0].VoteExtension, &voteExt); err != nil {
		h.logger.Error("failed to decode vote extension", "err", err, "validator", fmt.Sprintf("%x", ci.Votes[0].Validator.Address))
		return IDs, err
	}

	if len(voteExt.IDs) == 0 {
		return IDs, errors.New("no IDs in vote extension")
	}

	return voteExt.IDs, nil
}
