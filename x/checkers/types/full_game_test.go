package types_test

import (
	"strings"
	"testing"

	"github.com/bekauz/checkers/x/checkers/rules"
	"github.com/bekauz/checkers/x/checkers/testutil"
	"github.com/bekauz/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	alice = testutil.Alice
	bob   = testutil.Bob
)

func GetStoredGame1() types.StoredGame {
	return types.StoredGame{
		Black: alice,
		Red:   bob,
		Index: "1",
		Board: rules.New().String(),
		Turn:  "b",
	}
}

func TestCanGetAddressBlack(t *testing.T) {
	aliceAddress, err1 := sdk.AccAddressFromBech32(alice)
	black, err2 := GetStoredGame1().GetBlackAddress()
	require.Equal(t, aliceAddress, black)
	require.Nil(t, err2)
	require.Nil(t, err1)
}

func TestGetWrongAddressBlack(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Black = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d4"
	black, err := storedGame.GetBlackAddress()
	require.Nil(t, black)
	require.EqualError(t,
		err,
		"black address is invalid: cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d4: decoding bech32 failed: invalid checksum (expected 3xn9d3 got 3xn9d4)")
	require.EqualError(t, storedGame.Validate(), err.Error())
}

func TestCanGetAddressRed(t *testing.T) {
	bobAddress, err1 := sdk.AccAddressFromBech32(bob)
	red, err2 := GetStoredGame1().GetRedAddress()
	require.Equal(t, bobAddress, red)
	require.Nil(t, err1)
	require.Nil(t, err2)
}

func TestGetWrongAddressRed(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Red = "cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8h"
	red, err := storedGame.GetRedAddress()
	require.Nil(t, red)
	require.EqualError(t,
		err, "red address is invalid: cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8h: decoding bech32 failed: invalid checksum (expected xqhc8g got xqhc8h)")
	require.EqualError(t, storedGame.Validate(), err.Error())
}

func TestParseGame(t *testing.T) {
	storedGame := GetStoredGame1()
	game, err := storedGame.ParseGame()
	require.EqualValues(t, rules.New().Pieces, game.Pieces)
	require.Nil(t, err)
}

func TestParseErrorGame(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Board = "ha"
	game, err := storedGame.ParseGame()
	require.EqualError(t, err, "game is not parseable: invalid board string: ha")
	require.Nil(t, game)
}

func TestCheckCorrectTampering(t *testing.T) {
	storedGame := GetStoredGame1()
	// replace all black with red, technically correct
	storedGame.Board = strings.Replace(storedGame.Board, "b", "r", 1)
	game, err := storedGame.ParseGame()
	require.NotEqualValues(t, rules.New().Pieces, game.Pieces)
	require.Nil(t, err)
}

func TestCheckTampering(t *testing.T) {
	storedGame := GetStoredGame1()
	// replace black with green
	storedGame.Board = strings.Replace(storedGame.Board, "b", "g", 1)
	game, err := storedGame.ParseGame()
	require.Nil(t, game)
	require.EqualError(t, err, "game is not parseable: invalid board, invalid piece at 1, 0")
	require.EqualError(t, storedGame.Validate(), err.Error())
}
