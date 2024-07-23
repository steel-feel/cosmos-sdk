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
	"xarchain/x/xarchain/types"
)

// StakeWeightedPrices defines the structure a proposer should use to calculate
// and submit the stake-weighted prices for a given set of supported currency
// pairs, in addition to the vote extensions used to calculate them. This is so
// validators can verify the proposer's calculations.

// Its kind of temp storage in block formation
// before using this struct to commit the values to storage
type SuccessTransactionsID struct {
	IntentIDs          []uint64
	TxHashs            []string
	NextBlockHeight    int64
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

		if req.Height >= ctx.ConsensusParams().Abci.VoteExtensionsEnableHeight {
			/// NOTE: should not be commented out in production
			err := baseapp.ValidateVoteExtensions(ctx, h.valStore, req.Height, ctx.ChainID(), req.LocalLastCommit)
			if err != nil {
				return nil, err
			}

			cResp, err := h.computeCAIds(req.LocalLastCommit)
			if err != nil {
				return &abci.ResponsePrepareProposal{
					Txs: proposalTxs,
				}, nil
			}

			// if len(cResp.ComputedIDs) == 0 {
			// 	return nil, errors.New("no good transactions")
			// }

			injectedVoteExtTx := SuccessTransactionsID{
				IntentIDs:          cResp.ComputedIDs,
				TxHashs:            cResp.ComputedTxHashs,
				NextBlockHeight:    cResp.NextBlockHeight,
				ExtendedCommitInfo: req.LocalLastCommit,
			}

			h.logger.Warn("inside prepare proposal", "NextBlockHeight", cResp.NextBlockHeight)

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

		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}
}

func (h *ProposalHandler) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	res := &sdk.ResponsePreBlock{}
	h.logger.Warn("Txs len", "len tx", len(req.Txs))
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

		h.logger.Warn("inside Pre Blocker if condition")
	
		cBlock1 := types.Cblock{
			Blocknumber: injectedVoteExtTx.NextBlockHeight,
		}

		h.keeper.SetCblock(ctx, cBlock1)

		for i := 0; i < len(injectedVoteExtTx.IntentIDs); i++ {
			intent, found := h.keeper.GetIntentById(ctx, injectedVoteExtTx.IntentIDs[i])
			if !found {
				h.logger.Error("failed to get Intent %v", injectedVoteExtTx.IntentIDs[i])
				return nil, fmt.Errorf("failed to get Intent %v", injectedVoteExtTx.IntentIDs[i])
			}
			intent.Status = "verified"
			intent.Txhash = injectedVoteExtTx.TxHashs[i]
			//NOTE: filer address could also be fetched from vote extension
			//intent.Filer = fmt.Sprintf("%x", req.Validator.Address)

			h.keeper.SetIntent(ctx, intent)
		}

	}

	return res, nil
}

type ComputedResp struct {
	ComputedIDs     []uint64
	ComputedTxHashs []string
	NextBlockHeight int64
}

func (h *ProposalHandler) computeCAIds(ci abci.ExtendedCommitInfo) (*ComputedResp, error) {
	var voteExt CAVoteExtension
	if len(ci.Votes) == 0 {
		return nil, errors.New("no votes in commit info")
	}

	if err := json.Unmarshal(ci.Votes[0].VoteExtension, &voteExt); err != nil {
		h.logger.Error("failed to decode vote extension", "err", err, "validator", fmt.Sprintf("%x", ci.Votes[0].Validator.Address))
		return nil, err
	}

	// if len(voteExt.IDs) == 0 {
	// 	return nil, errors.New("no IDs in vote extension")
	// }

	return &ComputedResp{
		ComputedIDs:     voteExt.IDs,
		ComputedTxHashs: voteExt.TxHashs,
		NextBlockHeight: voteExt.Blocknumber,
	}, nil
}
