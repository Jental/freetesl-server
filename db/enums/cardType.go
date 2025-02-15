package enums

type CardType byte

const (
	CardTypeAction   CardType = 1
	CardTypeCreature CardType = 2
	CardTypeItem     CardType = 3
	CartTypeSupport  CardType = 4
)

var CardTypeName = map[CardType]string{
	CardTypeAction:   "Action",
	CardTypeCreature: "Creature",
	CardTypeItem:     "Item",
	CartTypeSupport:  "Support",
}
