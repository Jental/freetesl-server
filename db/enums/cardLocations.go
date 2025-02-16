package enums

type CardLocation byte

const (
	CardLocationDeck        CardLocation = 0
	CardLocationHand        CardLocation = 1
	CardLocationLeftLane    CardLocation = 2
	CardLocationRightLane   CardLocation = 3
	CardLocationDiscardPile CardLocation = 4
)
