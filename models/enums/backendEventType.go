package enums

type BackendEventType int

const (
	BackendEventDeckChanged                  BackendEventType = 0
	BackendEventHandChanged                  BackendEventType = 1
	BackendEventDiscardPileChanged           BackendEventType = 2
	BackendEventLanesChanged                 BackendEventType = 3
	BackendEventManaChanged                  BackendEventType = 4
	BackendEventHealthChanged                BackendEventType = 5
	BackendEventMatchStateRefresh            BackendEventType = 6
	BackendEventCardInstancesChanged         BackendEventType = 7
	BackendEventOpponentDeckChanged          BackendEventType = 100
	BackendEventOpponentHandChanged          BackendEventType = 101
	BackendEventOpponentDiscardPileChanged   BackendEventType = 102
	BackendEventOpponentLanesChanged         BackendEventType = 103
	BackendEventOpponentManaChanged          BackendEventType = 104
	BackendEventOpponentHealthChanged        BackendEventType = 105
	BackendEventOpponentMatchStateRefresh    BackendEventType = 106
	BackendEventOpponentCardInstancesChanged BackendEventType = 107
	BackendEventMatchStart                   BackendEventType = 200
	BackendEventMatchEnd                     BackendEventType = 201
	BackendEventSwitchTurn                   BackendEventType = 202
)

var BackendEventTypeName = map[BackendEventType]string{
	BackendEventDeckChanged:                  "DeckChanged",
	BackendEventHandChanged:                  "HandChanged",
	BackendEventDiscardPileChanged:           "DiscardPileChanged",
	BackendEventLanesChanged:                 "LanesChanged",
	BackendEventManaChanged:                  "ManaChanged",
	BackendEventHealthChanged:                "HealthChanged",
	BackendEventMatchStateRefresh:            "MatchStateRefresh",
	BackendEventCardInstancesChanged:         "CardInstancesChanged",
	BackendEventOpponentDeckChanged:          "OpponentDeckChanged",
	BackendEventOpponentHandChanged:          "OpponentHandChanged",
	BackendEventOpponentDiscardPileChanged:   "OpponentDiscardPileChanged",
	BackendEventOpponentLanesChanged:         "OpponentLanesChanged",
	BackendEventOpponentManaChanged:          "OpponentManaChanged",
	BackendEventOpponentHealthChanged:        "OpponentHealthChanged",
	BackendEventOpponentMatchStateRefresh:    "OpponentMatchStateRefresh",
	BackendEventOpponentCardInstancesChanged: "OpponentCardInstancesChanged",
	BackendEventMatchEnd:                     "MatchEnd",
	BackendEventSwitchTurn:                   "SwitchTurn",
}
