package keeper

import (
	"github.com/bekauz/checkers/x/checkers/types"
)

var _ types.QueryServer = Keeper{}
