package game

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var unshuffled = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 301, 302, 303, 304, 305, 306, 307, 308, 309, 310, 311, 312, 313}

func TestNewDeck(t *testing.T) {
	shuffled := NewDeck()
	require.ElementsMatch(t, unshuffled, shuffled)
	require.NotEqual(t, unshuffled, shuffled)
}

func TestShuffle(t *testing.T) {
	deck := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	shuffled := shuffle(deck)
	require.ElementsMatch(t, deck, shuffled)
	require.NotEqual(t, deck, shuffled)
}

func TestDraw(t *testing.T) {
	var hand []int
	deck := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	hand, deck = draw(-1, deck)
	require.Empty(t, hand)
	require.Len(t, deck, 10)

	hand, deck = draw(0, deck)
	require.Empty(t, hand)
	require.Len(t, deck, 10)

	hand, deck = draw(1, deck)
	require.Equal(t, []int{1}, hand)
	require.Len(t, deck, 9)

	hand, deck = draw(5, deck)
	require.Equal(t, []int{2, 3, 4, 5, 6}, hand)
	require.Len(t, deck, 4)

	hand, deck = draw(5, deck)
	require.Equal(t, []int{7, 8, 9, 10}, hand)
	require.Len(t, deck, 0)
}
