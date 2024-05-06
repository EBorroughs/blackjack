package game

import (
	"errors"
	"strconv"

	"blackjack/storage"
)

const (
	undecided = iota - 1
	loss
	win
	tie
)

var suits = []string{"♠", "♦", "♣", "♥"}

type State struct {
	deck       []int
	dealerHand []int
	playerHand []int
	wins       int
	losses     int
	ties       int
	result     int
}

func New(gameStateStorage storage.GameStateStorage, sessionID string) (*State, error) {
	s, err := gameStateStorage.GetGameState(sessionID)
	if err != nil {
		return nil, err
	}

	state := fromStorageGameState(s)
	if state == nil {
		state = &State{}
	}

	var dealerHand, playerHand []int
	deck := NewDeck()
	dealerHand, deck = draw(2, deck)
	playerHand, deck = draw(2, deck)
	state.deck = deck
	state.dealerHand = dealerHand
	state.playerHand = playerHand
	state.result = undecided
	gameStateStorage.UpsertGameState(sessionID, state.ToStorageGameState())
	return state, nil
}

func Get(gameStateStorage storage.GameStateStorage, sessionID string) (*State, error) {
	s, err := gameStateStorage.GetGameState(sessionID)
	if err != nil {
		return nil, err
	}

	return fromStorageGameState(s), nil
}

func Hit(gameStateStorage storage.GameStateStorage, sessionID string) (*State, error) {
	s, err := gameStateStorage.GetGameState(sessionID)
	if s == nil || err != nil {
		return nil, err
	}

	state := fromStorageGameState(s)
	if state.result != undecided {
		return state, errors.New("this game is already over, start a new game or reset the session to continue playing")
	}

	new, deck := draw(1, state.deck)
	state.deck = deck
	state.playerHand = append(state.playerHand, new...)
	if scoreHand(state.playerHand) > 21 {
		state.result = loss
		state.losses++
	}

	gameStateStorage.UpsertGameState(sessionID, state.ToStorageGameState())
	return state, nil
}

func Stand(gameStateStorage storage.GameStateStorage, sessionID string) (*State, error) {
	s, err := gameStateStorage.GetGameState(sessionID)
	if s == nil || err != nil {
		return nil, err
	}

	state := fromStorageGameState(s)
	if state.result != undecided {
		return state, errors.New("this game is already over, start a new game or reset the session to continue playing")
	}

	dealerScore := scoreHand(state.dealerHand)
	for dealerScore < 17 {
		new, deck := draw(1, state.deck)
		state.deck = deck
		state.dealerHand = append(state.dealerHand, new...)
		dealerScore = scoreHand(state.dealerHand)
	}

	playerScore := scoreHand(state.playerHand)
	if playerScore > dealerScore || dealerScore > 21 {
		state.result = win
		state.wins++
	} else if playerScore < dealerScore {
		state.result = loss
		state.losses++
	} else {
		state.result = tie
		state.ties++
	}

	gameStateStorage.UpsertGameState(sessionID, state.ToStorageGameState())
	return state, nil
}

func Delete(gameStateStorage storage.GameStateStorage, sessionID string) {
	gameStateStorage.DeleteGameState(sessionID)
}

func scoreHand(hand []int) int {
	var score, numAces int
	for _, card := range hand {
		switch value := card % 100; value {
		case 1:
			score += 11
			numAces++
		case 11, 12, 13:
			score += 10
		default:
			score += value
		}
	}

	for score > 21 && numAces > 0 {
		score -= 10
		numAces--
	}

	return score
}

func (s State) Wins() int {
	return s.wins
}

func (s State) Losses() int {
	return s.losses
}

func (s State) Ties() int {
	return s.ties
}

func (s State) Result() int {
	return s.result
}

func (s State) ToCardFormat() ([]string, []string) {
	var dealerHand, playerHand []string
	dealerHand = append(dealerHand, intToCard(s.dealerHand[0]))
	for i := 1; i < len(s.dealerHand); i++ {
		if s.result != undecided {
			dealerHand = append(dealerHand, intToCard(s.dealerHand[i]))
			continue
		}
		dealerHand = append(dealerHand, "?")
	}

	for _, card := range s.playerHand {
		playerHand = append(playerHand, intToCard(card))
	}

	return dealerHand, playerHand
}

func intToCard(cardInt int) string {
	var card string
	switch value := cardInt % 100; value {
	case 1:
		card = "A"
	case 11:
		card = "J"
	case 12:
		card = "Q"
	case 13:
		card = "K"
	default:
		card = strconv.Itoa(value)
	}

	return card + suits[cardInt/100]
}

func (s State) ToStorageGameState() *storage.GameState {
	return &storage.GameState{
		Deck:       s.deck,
		DealerHand: s.dealerHand,
		PlayerHand: s.playerHand,
		Wins:       s.wins,
		Losses:     s.losses,
		Ties:       s.ties,
		Result:     s.result,
	}
}

func fromStorageGameState(gameState *storage.GameState) *State {
	if gameState == nil {
		return nil
	}

	return &State{
		deck:       gameState.Deck,
		dealerHand: gameState.DealerHand,
		playerHand: gameState.PlayerHand,
		wins:       gameState.Wins,
		losses:     gameState.Losses,
		ties:       gameState.Ties,
		result:     gameState.Result,
	}
}
