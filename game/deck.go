package game

import (
	"math/rand/v2"
)

func NewDeck() []int {
	deck := make([]int, 0, 52)
	for i := 0; i < 4; i++ {
		for j := 1; j < 14; j++ {
			deck = append(deck, 100*i+j)
		}
	}

	return shuffle(deck)
}

func shuffle(deck []int) []int {
	tmp := make([]int, len(deck))
	copy(tmp, deck)
	rand.Shuffle(len(tmp), func(i, j int) {
		tmp[i], tmp[j] = tmp[j], tmp[i]
	})

	return tmp
}

func draw(numCards int, deck []int) ([]int, []int) {
	if numCards <= 0 {
		return nil, deck
	}

	tmp := make([]int, len(deck))
	copy(tmp, deck)
	if numCards > len(tmp) {
		return tmp, nil
	}

	return tmp[0:numCards], tmp[numCards:]
}
