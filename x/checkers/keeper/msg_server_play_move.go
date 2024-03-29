package keeper

import (
	"context"
	"strconv"

	"github.com/bekauz/checkers/x/checkers/rules"
	"github.com/bekauz/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) PlayMove(goCtx context.Context, msg *types.MsgPlayMove) (*types.MsgPlayMoveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get the stored game
	storedGame, found := k.Keeper.GetStoredGame(ctx, msg.GameIndex)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "%s", msg.GameIndex)
	}

	// determine player color
	isBlack := storedGame.Black == msg.Creator
	isRed := storedGame.Red == msg.Creator
	var player rules.Player
	if !isBlack && isRed {
		player = rules.RED_PLAYER
	} else if isBlack && !isRed {
		player = rules.BLACK_PLAYER
	} else {
		player = rules.StringPieces[storedGame.Turn].Player
	}
	// parse the game
	game, err := storedGame.ParseGame()
	if err != nil {
		panic(err.Error())
	}

	// validate the player turn
	if !game.TurnIs(player) {
		return nil, sdkerrors.Wrapf(types.ErrNotPlayerTurn, "%s", player)
	}

	captured, moveErr := game.Move(
		rules.Pos{
			X: int(msg.FromX),
			Y: int(msg.FromY),
		},
		rules.Pos{
			X: int(msg.ToX),
			Y: int(msg.ToY),
		},
	)
	if moveErr != nil {
		return nil, sdkerrors.Wrapf(types.ErrWrongMove, moveErr.Error())
	}

	storedGame.Board = game.String()
	storedGame.Turn = rules.PieceStrings[game.Turn]
	k.Keeper.SetStoredGame(ctx, storedGame)

	// emit the move event
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.MovePlayedEventType,
		sdk.NewAttribute(types.MovePlayedEventCapturedX, strconv.FormatInt(int64(captured.X), 10)),
		sdk.NewAttribute(types.MovePlayedEventCapturedY, strconv.FormatInt(int64(captured.Y), 10)),
		sdk.NewAttribute(types.MovePlayedEventCreator, msg.Creator),
		sdk.NewAttribute(types.MovePlayedEventGameIndex, msg.GameIndex),
		sdk.NewAttribute(types.MovePlayedEventWinner, rules.PieceStrings[game.Winner()]),
	))

	return &types.MsgPlayMoveResponse{
		CapturedX: int32(captured.X),
		CapturedY: int32(captured.Y),
		Winner:    rules.PieceStrings[game.Winner()],
	}, nil
}
