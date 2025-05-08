package main

import (
	"strings"

	"github.com/whywaita/poker-go"
)

type Card struct {
	c poker.Card
}

func PokerCardString(pc poker.Card) string {
	if IsUnknownCard(pc) {
		return "unknown"
	}
	sb := strings.Builder{}
	// https://gist.github.com/JBlond/2fea43a3049b38287e5e9cefc87b2124
	if pc.Suit == poker.Hearts {
		sb.WriteString("\x1b[31m")
	}
	if pc.Suit == poker.Diamonds {
		sb.WriteString("\x1b[34m")
	}
	if pc.Suit == poker.Spades {
		sb.WriteString("\x1b[38;5;238m")
	}
	if pc.Suit == poker.Clubs {
		sb.WriteString("\x1b[32m")
	}
	sb.WriteString(pc.Rank.String())
	sb.WriteString(SuitStringToUnicode(pc.Suit.String()))
	sb.WriteString("\x1b[0m")
	return sb.String()
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

