package main

import (
	"github.com/whywaita/poker-go"
)

type Card struct {
	c poker.Card
}

func PokerCardString(pc poker.Card) string {
	if IsUnknownCard(pc) {
		return "unknown"
	}
	return pc.Rank.String() + "" + SuitStringToUnicode(pc.Suit.String())
}

func NewCard(c poker.Card) Card {
	return Card{c: c}
}

func (c Card) String() string {
	pc := poker.Card(c.c)
	return PokerCardString(pc)
}

func SuitStringToUnicode(suit string) string {
	switch suit {
	case "hearts":
		return "♥"
	case "diamonds":
		return "♦"
	case "clubs":
		return "♣"
	case "spades":
		return "♠"
	}

	return "[Unknown]"
}

func (c Card) Id() string {
	pc := poker.Card(c.c)
	return pc.Rank.String() + "_" + pc.Suit.String()
}

func (c Card) Title() string       { return c.String() }
func (c Card) Description() string { return c.String() }

func (c Card) FilterValue() string { return c.String() }
