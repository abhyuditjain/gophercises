package deck

import (
	"fmt"
	"math/rand"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Spade})
	fmt.Println(Card{Rank: Two, Suit: Diamond})
	fmt.Println(Card{Rank: Queen, Suit: Heart})
	fmt.Println(Card{Rank: King, Suit: Club})
	fmt.Println(Card{Rank: Jack, Suit: Joker})

	// Output:
	// Ace of Spades
	// Two of Diamonds
	// Queen of Hearts
	// King of Clubs
	// Joker
}

func TestNew(t *testing.T) {
	deck := New()

	if len(deck) != 52 {
		t.Errorf("Wrong number of cards in a new deck: want %d, got %d", 52, len(deck))
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	want := Card{Suit: Spade, Rank: Ace}
	if cards[0] != want {
		t.Errorf("First card: got %+v, want %+v", cards[0], want)
	}
	want = Card{Suit: Heart, Rank: King}
	if cards[len(cards)-1] != want {
		t.Errorf("Last card: got %+v, want %+v", cards[len(cards)-1], want)
	}
}

func TestSort(t *testing.T) {
	cards := New(Sort(Less))
	want := Card{Suit: Spade, Rank: Ace}
	if cards[0] != want {
		t.Errorf("First card: got %+v, want %+v", cards[0], want)
	}
	want = Card{Suit: Heart, Rank: King}
	if cards[len(cards)-1] != want {
		t.Errorf("Last card: got %+v, want %+v", cards[len(cards)-1], want)
	}
}

func TestJokers(t *testing.T) {
	cards := New(Jokers(3))
	count := 0
	for _, card := range cards {
		if card.Suit == Joker {
			count++
		}
	}
	if count != 3 {
		t.Errorf("jokers count got %d, want %d", count, 3)
	}
}

func TestFilter(t *testing.T) {
	// Filter function to return only spades
	filter := func(c Card) bool {
		return c.Suit == Spade
	}

	cards := New(Filter(filter))
	for _, v := range cards {
		if v.Suit != Spade {
			t.Errorf("Expected only spades, but got suit: %v", v.Suit)
		}
	}
}

func TestDeck(t *testing.T) {
	cards := New(Deck(3))
	if len(cards) != 3*52 {
		t.Errorf("Wanted %d decks, got %d decks", 3, len(cards)/52)
	}
}

func TestShuffle(t *testing.T) {
	// make shuffleRand deterministic
	// First call to shuffleRand.Perm(52) should be:
	// [40, 35 ...]
	shuffleRand = rand.New(rand.NewSource(0))

	orig, shuffled := New(), New(Shuffle)

	if shuffled[0] != orig[40] {
		t.Errorf("Expected first card in shuffled deck to be %q, but got %q", orig[40], shuffled[0])
	}
	if shuffled[1] != orig[35] {
		t.Errorf("Expected second card in shuffled deck to be %q, but got %q", orig[35], shuffled[1])
	}
}
