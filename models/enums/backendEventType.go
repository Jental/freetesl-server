package enums

type BackendEventType int

const (
	BackendEventDeckChanged                BackendEventType = 0
	BackendEventHandChanged                BackendEventType = 1
	BackendEventDiscardPileChanged         BackendEventType = 2
	BackendEventLanesChanged               BackendEventType = 3
	BackendEventManaChanged                BackendEventType = 4
	BackendEventHealthChanged              BackendEventType = 5
	BackendEventMatchStateRefresh          BackendEventType = 6
	BackendEventOpponentDeckChanged        BackendEventType = 100
	BackendEventOpponentHandChanged        BackendEventType = 101
	BackendEventOpponentDiscardPileChanged BackendEventType = 102
	BackendEventOpponentLanesChanged       BackendEventType = 103
	BackendEventOpponentManaChanged        BackendEventType = 104
	BackendEventOpponentHealthChanged      BackendEventType = 105
	BackendEventOpponentMatchStateRefresh  BackendEventType = 106
	BackendEventMatchEnd                   BackendEventType = 200
	BackendEventSwitchTurn                 BackendEventType = 201
)

var BackendEventTypeName = map[BackendEventType]string{
	BackendEventDeckChanged:                "DeckChanged",
	BackendEventHandChanged:                "HandChanged",
	BackendEventDiscardPileChanged:         "DiscardPileChanged",
	BackendEventLanesChanged:               "LanesChanged",
	BackendEventManaChanged:                "ManaChanged",
	BackendEventHealthChanged:              "HealthChanged",
	BackendEventMatchStateRefresh:          "MatchStateRefresh",
	BackendEventOpponentDeckChanged:        "OpponentDeckChanged",
	BackendEventOpponentHandChanged:        "OpponentHandChanged",
	BackendEventOpponentDiscardPileChanged: "OpponentDiscardPileChanged",
	BackendEventOpponentLanesChanged:       "OpponentLanesChanged",
	BackendEventOpponentManaChanged:        "OpponentManaChanged",
	BackendEventOpponentHealthChanged:      "OpponentHealthChanged",
	BackendEventOpponentMatchStateRefresh:  "OpponentMatchStateRefresh",
	BackendEventMatchEnd:                   "MatchEnd",
	BackendEventSwitchTurn:                 "SwitchTurn",
}
