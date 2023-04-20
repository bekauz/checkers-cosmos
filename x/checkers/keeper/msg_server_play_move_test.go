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

func setupMsgServerWithOneGameForPlayMove(t testing.TB) (types.MsgServer, keeper.Keeper, context.Context) {
	k, ctx := keepertest.CheckersKeeper(t)
	checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
	server := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)
	server.CreateGame(context, &types.MsgCreateGame{
		Creator: testutil.Alice,
		Black:   testutil.Bob,
		Red:     testutil.Carol,
	})
	return server, *k, context
}

func TestPlayMoveHappy(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForPlayMove(t)

	playMoveResponse, err := msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   testutil.Bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})

	require.Nil(t, err)
	require.EqualValues(t, &types.MsgPlayMoveResponse{
		CapturedX: -1,
		CapturedY: -1,
		Winner:    "*",
	}, playMoveResponse)
}

func TestPlayMoveNoGame(t *testing.T) {
	k, ctx := keepertest.CheckersKeeper(t)
	checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
	msgServer := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)

	_, err := msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   testutil.Bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})

	require.NotNil(t, err)
	require.Equal(t, "1: game not found", err.Error())
}

func TestPlayMoveNotPlayerTurn(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForPlayMove(t)

	_, err := msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   testutil.Carol,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})

	require.NotNil(t, err)
	require.Equal(t, "{red}: player tried to play out of turn", err.Error())
}

func TestPlayMoveInvalidGame(t *testing.T) {
	// TODO
}

func TestPlayMoveWrongMove(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForPlayMove(t)

	_, err := msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   testutil.Bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       55,
	})

	require.NotNil(t, err)
	require.Equal(t, "Invalid move: {1 2} to {2 55}: wrong move", err.Error())
}
