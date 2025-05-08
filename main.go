package main

import (
	"fmt"

	"github.com/whywaita/poker-go"
)

var allCards [52]poker.Card

func init() {
	newDeck := poker.NewDeck()
	for i := range 52 {
		allCards[i] = newDeck.DrawCard()
	}
}

func main() {

	for _, c := range allCards {
		fmt.Println(c)
	}
}
