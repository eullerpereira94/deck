//go:generate stringer -type=Suit,Rank

package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// Suit is an abstraction of a card s
type Suit uint8

// The current possible suits, plus the special case of a Joker and the arcanes of Tarot
const (
	// Spade suit
	Spade Suit = iota

	// Diamond suit
	Diamond

	// Club suit
	Club

	// Hear suit
	Heart

	// Joker is a special case
	Joker
)

var suits = [...]Suit{Spade, Diamond, Club, Heart}

// Rank is an abstraction of possible ranks a card might have
type Rank uint8

// The possible common ranks for playing deck plus a Knight rank between the Jack and the Queen
const (
	_ Rank = iota

	// Ace rank
	Ace

	// Two rank
	Two

	// Three rank
	Three

	// Four rank
	Four

	// Five rank
	Five

	// Six rank
	Six

	// Seven rank
	Seven

	// Eight rank
	Eight

	// Nine rank
	Nine

	// Ten rank
	Ten

	// Jack rank
	Jack

	// Queen rank
	Queen

	// King rank
	King
)

const (
	minRank = Ace
	maxRank = King
)

// Card is the logical representation of a single card in a traditional playing deck
type Card struct {
	Suit
	Rank
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}

	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

// New returns a deck of cards with the specified number of joker and with an option to include knights
func New(opts ...func([]Card) []Card) []Card {
	var cards []Card

	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}

	for _, opt := range opts {
		cards = opt(cards)
	}

	return cards
}

// DefaultSort is the default sorting of a newly created deck, first by rank, then by suit.
func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, BySuitThenByRank(cards))

	return cards
}

// Sort is a funtion that receives a comparator of cards and returns a function that sorts a deck based on said comparation.
func Sort(less func(cards []Card) func(i, j int) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

// ByRankThenBySuit receives a slice of Card and returns a function that compares two cards first by rank, then by suit.
func ByRankThenBySuit(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return rankThenSuit(cards[i]) < rankThenSuit(cards[j])
	}
}

// BySuitThenByRank creates a function that helps to fullfil the contract defined by sort.Interface interface.
func BySuitThenByRank(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank(cards[i]) < absRank(cards[j])
	}
}

func absRank(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

func rankThenSuit(c Card) int {
	return ((int(c.Rank) - 1) * len(suits)) + int(c.Suit)
}

// Shuffle  distribute the elements of the slice in a random order
func Shuffle(cards []Card) []Card {
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	r.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	return cards
}

// Jokers inserts n Joker cards in our deck
func Jokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		for i := 0; i < n; i++ {
			cards = append(cards, Card{Suit: Joker, Rank: Rank(i)})
		}
		return cards
	}
}

// Filter takes a filter function and returns a function that receives a slice of cards and returns a slice with all the elements that weren't excluded by the filter
func Filter(f func(card Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		ret := []Card{}

		for _, card := range cards {
			if !f(card) {
				ret = append(ret, card)
			}
		}

		return ret
	}
}

// Deck creates n slices of Card that are identical copies of the deck that would be generated by other option functions
func Deck(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for i := 0; i < n; i++ {
			ret = append(ret, cards...)
		}
		return ret
	}
}
