package main

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strings"
	"time"
)

const (
	MAIN_STARTING_HP = 20
)

type Deck []Card

func (d *Deck) isEmpty() bool {
	return len(MAIN_DECK) == 0
}

// func (d *Deck) hasLastCard() bool {
// 	return len(MAIN_DECK) == 1
// }

var (
	MAIN_DECK      = Deck{}
	MAIN_ARENA     = [4]Card{} // only upto 4 cards in the arena are allowed
	CARDS_IN_ARENA = 0
	DISCARDED_DECK = Deck{} // i dont this we need this
	MAIN_CURRENT   = Current{}
	SUITS          = []CardSuit{
		CardSuitClubs,
		CardSuitSpades,
		CardSuitHearts,
		CardSuitDiamonds,
	}
	SYMBOLS = []CardSymbol{
		CardSymbolOne,
		CardSymbolTwo,
		CardSymbolThree,
		CardSymbolFour,
		CardSymbolFive,
		CardSymbolSix,
		CardSymbolSeven,
		CardSymbolEight,
		CardSymbolNine,
		CardSymbolTen,

		CardSymbolJack,
		CardSymbolQueen,
		CardSymbolKing,
		CardSymbolAce,
	}
	MAIN_HP         = MAIN_STARTING_HP // starts with 20
	AvoidedLastRoom = false
)

type CardSuit string

const (
	CardSuitHearts   = "♥"
	CardSuitClubs    = "♣"
	CardSuitSpades   = "♠"
	CardSuitDiamonds = "♦"
)

type CardSymbol string

const (
	CardSymbolOne   CardSymbol = "1"
	CardSymbolTwo   CardSymbol = "2"
	CardSymbolThree CardSymbol = "3"
	CardSymbolFour  CardSymbol = "4"
	CardSymbolFive  CardSymbol = "5"
	CardSymbolSix   CardSymbol = "6"
	CardSymbolSeven CardSymbol = "7"
	CardSymbolEight CardSymbol = "8"
	CardSymbolNine  CardSymbol = "9"
	CardSymbolTen   CardSymbol = "10"

	CardSymbolJack  CardSymbol = "J"
	CardSymbolQueen CardSymbol = "Q"
	CardSymbolKing  CardSymbol = "K"
	CardSymbolAce   CardSymbol = "A"
)

type CardType string

const (
	CardTypeMonster = "monster"
	CardTypeWeapon  = "weapon"
	CardTypeHP      = "hp"
)

type Card struct {
	Suit       CardSuit
	Symbol     CardSymbol
	IsFaceCard bool
	Type       CardType // monster, weapon or hp
	Value      int      // 1-10 and 11, 12, 13, 14 for face cards and aces respectively
}

type Current struct {
	Weapon       *Card
	LastDefeated *Card
}

func getCardTypeBySuit(suit CardSuit) (t CardType) {
	switch suit {
	case CardSuitClubs, CardSuitSpades:
		t = CardTypeMonster
	case CardSuitDiamonds:
		t = CardTypeWeapon
	case CardSuitHearts:
		t = CardTypeHP
	}
	return
}

func getCardValueBySymbol(symbol CardSymbol) (v int) {
	switch symbol {

	case CardSymbolOne:
		v = 1
	case CardSymbolTwo:
		v = 2
	case CardSymbolThree:
		v = 3
	case CardSymbolFour:
		v = 4
	case CardSymbolFive:
		v = 5
	case CardSymbolSix:
		v = 6
	case CardSymbolSeven:
		v = 7
	case CardSymbolEight:
		v = 8
	case CardSymbolNine:
		v = 9
	case CardSymbolTen:
		v = 10
	case CardSymbolJack:
		v = 11
	case CardSymbolQueen:
		v = 12
	case CardSymbolKing:
		v = 13
	case CardSymbolAce:
		v = 14
	}
	return
}

func removeCardFromArena(s int) {
	CARDS_IN_ARENA--
	a := append(MAIN_ARENA[:s], MAIN_ARENA[s+1:]...)
	copy(MAIN_ARENA[:], a)
}

func addCardToArena(c Card) {
	MAIN_ARENA[CARDS_IN_ARENA] = c
	CARDS_IN_ARENA++
}

func takeCardFromDeck() Card {
	c := MAIN_DECK[0]
	MAIN_DECK = MAIN_DECK[1:]
	return c
}

func clearMainArena() {
	CARDS_IN_ARENA = 0
}

func shuffleDeck() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(MAIN_DECK), func(i, j int) {
		MAIN_DECK[i], MAIN_DECK[j] = MAIN_DECK[j], MAIN_DECK[i]
	})
}

func initMainDeck() {
	// generate cards
	for _, suit := range SUITS {
		for _, symbol := range SYMBOLS {
			isFaceCard := slices.Contains([]CardSymbol{CardSymbolAce, CardSymbolKing, CardSymbolQueen, CardSymbolJack}, symbol)
			if (suit == CardSuitHearts || suit == CardSuitDiamonds) && isFaceCard {
				continue
			}
			card := Card{
				Suit:       suit,
				Symbol:     symbol,
				IsFaceCard: isFaceCard,
				Type:       getCardTypeBySuit(suit),
				Value:      getCardValueBySymbol(symbol),
			}
			MAIN_DECK = append(MAIN_DECK, card)
		}
	}

	// ensure shuffling
	shuffleDeck()
}

func printArena() {
	fmt.Println("-----ARENA-----")
	if CARDS_IN_ARENA == 0 {
		fmt.Println("No cards in arena yet.")
		return
	}

	for i, c := range MAIN_ARENA {
		fmt.Printf("%d. %s %s\n", i+1, c.Symbol, c.Suit)
	}
}

func newRoom() {
	// if the number of cards in the arena right now is 0 - extract 4 cards from the deck
	// if the number of cards in the arena is n (n > 0) - extract 4 - n cards from the deck
	for range 4 - CARDS_IN_ARENA {
		card := takeCardFromDeck()
		addCardToArena(card)
	}
}

func avoidRoom() {
	// move the cards in the arena (the room) - to the bottom of the deck
	MAIN_DECK = append(MAIN_DECK, MAIN_ARENA[0], MAIN_ARENA[1], MAIN_ARENA[2], MAIN_ARENA[3])
	clearMainArena()

	AvoidedLastRoom = true

	fmt.Println("avoided room")
}

func faceRoom() {
	AvoidedLastRoom = false
	// interact with 3 of the 4 cards in the room - ONE BY ONE
	for range 3 {
		var validInteractedCardNum bool
		var interactedCardNum int
		for !validInteractedCardNum {
			fmt.Print("Choose card to interact by number (1, 2, 3 or 4): ")
			fmt.Scanln(&interactedCardNum)
			if interactedCardNum <= CARDS_IN_ARENA && slices.Contains([]int{1, 2, 3, 4}, interactedCardNum) {
				validInteractedCardNum = true
			}
		}

		interactedCard := MAIN_ARENA[interactedCardNum-1]

		switch interactedCard.Type {
		case CardTypeHP:
			MAIN_HP = int(math.Min(MAIN_STARTING_HP, float64(MAIN_HP+interactedCard.Value)))
		case CardTypeMonster:
			damage := interactedCard.Value

			if MAIN_CURRENT.Weapon != nil {
				// ask the player if they want to use the equipped weapon
				var useEquippedWeaponStr string
				fmt.Print("Do you want to use the equipped weapon on this monster? (y/N): ")
				fmt.Scanln(&useEquippedWeaponStr)

				useEquippedWeapon := strings.ToLower(useEquippedWeaponStr) == "y"
				if useEquippedWeapon {
					if MAIN_CURRENT.LastDefeated != nil {
						if interactedCard.Value < MAIN_CURRENT.LastDefeated.Value {
							damage = int(math.Max(0, float64(interactedCard.Value-MAIN_CURRENT.Weapon.Value)))
							MAIN_CURRENT.LastDefeated = &interactedCard
						}
					}
				}
			}
			fmt.Printf("Faced monster with damage %d\n", damage)
			MAIN_HP = int(math.Max(0, float64(MAIN_HP-damage)))
			if MAIN_HP == 0 {
				return // game over
			}

		case CardTypeWeapon:
			MAIN_CURRENT = Current{
				Weapon: &interactedCard,
			}
		}
		removeCardFromArena(interactedCardNum - 1)
		fmt.Printf("HP: %d ❤️ \n", MAIN_HP)
	}
	// this means - each interaction changes the game state not the entire card facing
	fmt.Println("faced room")
}

func init() {
	initMainDeck() // ensures a new game is started with a new shuffled deck
}

func main() {
	gameLoop := true

	fmt.Println("Starting game...")
	for gameLoop && MAIN_HP > 0 {
		fmt.Printf("HP: %d ❤️ \n", MAIN_HP)
		fmt.Println("New Room...")
		newRoom()

		printArena()
		var validRoomChoice bool
		var roomChoice int
		for !validRoomChoice {
			if AvoidedLastRoom {
				fmt.Println("No choice, just Face (you avoided the last room)")
				roomChoice = 1
				validRoomChoice = true
			} else {
				fmt.Print("Choice (1. Face | 2. Avoid): ")
				fmt.Scanln(&roomChoice)
				if roomChoice == 1 || roomChoice == 2 {
					validRoomChoice = true
				}
				fmt.Println()
			}
		}

		if roomChoice == 1 {
			faceRoom()
		} else {
			avoidRoom()
		}

		if MAIN_DECK.isEmpty() {
			fmt.Println("Dungeon crossed! You win!")
			break
		} else if MAIN_HP == 0 {
			fmt.Println("You died. Game Over!")
		}
	}
}
