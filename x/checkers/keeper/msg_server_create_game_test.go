package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/bekauz/checkers/testutil/keeper"
	"github.com/bekauz/checkers/x/checkers"
	"github.com/bekauz/checkers/x/checkers/keeper"
	testutil "github.com/bekauz/checkers/x/checkers/testutil"
	"github.com/bekauz/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServerCreateGame(t testing.TB) (types.MsgServer, keeper.Keeper, context.Context) {
	k, ctx := keepertest.CheckersKeeper(t)
	checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
	return keeper.NewMsgServerImpl(*k), *k, sdk.WrapSDKContext(ctx)
}

func TestCreateGame(t *testing.T) {
	msgServer, _, context := setupMsgServerCreateGame(t)
	createResponse, err := msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: testutil.Alice,
		Black:   testutil.Bob,
		Red:     testutil.Carol,
	})
	require.Nil(t, err)
	require.EqualValues(t, types.MsgCreateGameResponse{
		GameIndex: "1",
	}, *createResponse)

	createResponse2, err2 := msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: testutil.Bob,
		Black:   testutil.Alice,
		Red:     testutil.Carol,
	})
	require.Nil(t, err2)
	require.EqualValues(t, types.MsgCreateGameResponse{
		GameIndex: "2",
	}, *createResponse2)
}
