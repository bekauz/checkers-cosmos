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

func TestCreateGameHasSaved(t *testing.T) {
	msgServer, keeper, context := setupMsgServerCreateGame(t)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: testutil.Alice,
		Black:   testutil.Bob,
		Red:     testutil.Carol,
	})
	systemInfo, found := keeper.GetSystemInfo(sdk.UnwrapSDKContext(context))
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId: 2,
	}, systemInfo)

	game, found := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index: "1",
		Board: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:  "b",
		Black: testutil.Bob,
		Red:   testutil.Carol,
	}, game)
}

func TestCreateMultipleGames(t *testing.T) {
	msgServer, keeper, context := setupMsgServerCreateGame(t)

	// create multiple games
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: testutil.Alice,
		Black:   testutil.Bob,
		Red:     testutil.Carol,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: testutil.Alice,
		Black:   testutil.Carol,
		Red:     testutil.Bob,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: testutil.Bob,
		Black:   testutil.Alice,
		Red:     testutil.Carol,
	})

	systemInfo, found := keeper.GetSystemInfo(sdk.UnwrapSDKContext(context))
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId: 4,
	}, systemInfo)

	game1, found := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index: "1",
		Board: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:  "b",
		Black: testutil.Bob,
		Red:   testutil.Carol,
	}, game1)

	game2, found := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "2")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index: "2",
		Board: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:  "b",
		Black: testutil.Carol,
		Red:   testutil.Bob,
	}, game2)

	game3, found := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "3")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index: "3",
		Board: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:  "b",
		Black: testutil.Alice,
		Red:   testutil.Carol,
	}, game3)
}

func TestCreate1GameEmitted(t *testing.T) {
	msgServer, _, context := setupMsgServerCreateGame(t)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: testutil.Alice,
		Black:   testutil.Bob,
		Red:     testutil.Carol,
	})
	ctx := sdk.UnwrapSDKContext(context)
	require.NotNil(t, ctx)
	// grab the resulting events after creating a game
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 1)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "new-game-created",
		Attributes: []sdk.Attribute{
			{Key: "creator", Value: testutil.Alice},
			{Key: "game-index", Value: "1"},
			{Key: "black", Value: testutil.Bob},
			{Key: "red", Value: testutil.Carol},
		},
	}, event)
}
