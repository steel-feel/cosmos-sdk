package abci

import (
	"encoding/json"
	"errors"
	"fmt"

	"cosmossdk.io/log"

	abci "github.com/cometbft/cometbft/abci/types"
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
type VoteExtensionTransaction struct {
	IntentData         map[string]IntentData // map of chainId to intent data
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
			/// NOTE: should not be commented out in production
			err := baseapp.ValidateVoteExtensions(ctx, h.valStore, req.Height, ctx.ChainID(), req.LocalLastCommit)
			if err != nil {
				return nil, err
			}

			intentData, err := h.computeCAIds(req.LocalLastCommit)
			if err != nil {
				return &abci.ResponsePrepareProposal{
					Txs: proposalTxs,
				}, nil
			}

			injectedVoteExtTx := VoteExtensionTransaction{
				IntentData:         intentData,
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
			var injectedVoteExtTx VoteExtensionTransaction
			if err := json.Unmarshal(req.Txs[0], &injectedVoteExtTx); err != nil {
				h.logger.Error("failed to decode proccess Proposal", "err", err)
				return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
			}

			// NOTE: We can validate extra things here vote extension here,
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
		var injectedVoteExtTx VoteExtensionTransaction
		if err := json.Unmarshal(req.Txs[0], &injectedVoteExtTx); err != nil {
			h.logger.Error("failed to decode preblocker vote", "err", err)
			// return nil, err
			return res, nil

		}

		//TODO: Upgrade using go routine
		for chainId, intentData := range injectedVoteExtTx.IntentData {
			// Set the intent data in the keeper
			if intentData.From != 0 && intentData.To != 0 {
			h.keeper.SetSyncblock(ctx, types.Syncblock{
				From:    intentData.From,
				To:      intentData.To,
				ChainId: chainId,
			})
			}

			for _, ID := range intentData.IDs {
				intent, found := h.keeper.GetIntentById(ctx, ID)

				if !found {
					h.logger.Error("failed to get Intent %v", ID)
					return nil, fmt.Errorf("failed to get Intent %v", ID)
				}
				intent.Status = "verified"
				h.keeper.SetIntent(ctx, intent)
			}
		}
	}
	return res, nil
}

type ComputedResp struct {
	ComputedIntents map[string]IntentData
}

func (h *ProposalHandler) computeCAIds(ci abci.ExtendedCommitInfo) (map[string]IntentData, error) {
	var voteExt CAVoteExtension
	if len(ci.Votes) == 0 {
		return nil, errors.New("no votes in commit info")
	}

	if err := json.Unmarshal(ci.Votes[0].VoteExtension, &voteExt); err != nil {
		h.logger.Error("failed to decode vote extension", "err", err, "validator", fmt.Sprintf("%x", ci.Votes[0].Validator.Address))
		return nil, err
	}

	return voteExt.EventData, nil
}
