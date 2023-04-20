package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/checkers module sentinel errors
var (
	ErrInvalidBlack         = sdkerrors.Register(ModuleName, 1100, "black address is invalid: %s")
	ErrInvalidRed           = sdkerrors.Register(ModuleName, 1101, "red address is invalid: %s")
	ErrGameNotParseable     = sdkerrors.Register(ModuleName, 1102, "game is not parseable")
	ErrInvalidGameIndex     = sdkerrors.Register(ModuleName, 1103, "game index is invalid")
	ErrInvalidPositionIndex = sdkerrors.Register(ModuleName, 1104, "position index is invalid")
	ErrMoveAbsent           = sdkerrors.Register(ModuleName, 1105, "move is absent")
)
