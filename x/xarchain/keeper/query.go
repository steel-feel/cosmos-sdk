package keeper

import (
	"xarchain/x/xarchain/types"
)

var _ types.QueryServer = Keeper{}
