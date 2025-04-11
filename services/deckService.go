package services

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/db/queries"
	"github.com/jental/freetesl-server/models"
	"github.com/samber/lo"
)

func GetDecks(playerID int) ([]*models.Deck, error) {
	decksFromDB, err := queries.GetDecks(playerID)
	if err != nil {
		return nil, err
	}

	allCards, err := GetAllCardsMap()
	if err != nil {
		return nil, err
	}

	allCardClasses, err := GetAllCardClasses()
	if err != nil {
		return nil, err
	}

	errs := make([]error, 0)

	decks := lo.Map(decksFromDB, func(dbDeck *dbModels.Deck, _ int) *models.Deck {
		deck := models.Deck{
			ID:         dbDeck.ID,
			Name:       dbDeck.Name,
			AvatarName: dbDeck.AvatarName,
			PlayerID:   dbDeck.PlayerID,
			Cards:      make([]*models.CardWithCount, 0),
		}

		attributes := make(map[int]*dbModels.Attribute)

		for cardID, count := range dbDeck.Cards {
			card, exists := allCards[cardID]
			if !exists {
				errs = append(errs, fmt.Errorf("deck [%d] have unexisting card: '%d'", deck.ID, cardID))
				continue
			}

			cardWithCount := models.CardWithCount{
				Card:  card,
				Count: count,
			}

			deck.Cards = append(deck.Cards, &cardWithCount)

			cardClass, exists := allCardClasses[card.ClassID]
			if !exists {
				errs = append(errs, fmt.Errorf("deck [%d] have a card [%d] with an unknown class: '%d'", deck.ID, cardID, card.ClassID))
			}
			for _, attribute := range cardClass.Attributes {
				_, exists := attributes[attribute.ID]
				if !exists {
					attributes[attribute.ID] = attribute
				}
			}
		}

		deck.Attributes = slices.Collect(maps.Values(attributes))

		return &deck
	})

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return decks, nil
}

func GetDeck(playerID int, deckID int) (*models.Deck, error) {
	decks, err := GetDecks(playerID) // TODO: find deck in db by id
	if err != nil {
		return nil, err
	}

	deckIdx := slices.IndexFunc(decks, func(d *models.Deck) bool { return d.ID == deckID })
	if deckIdx < 0 {
		return nil, fmt.Errorf("[%d]: deck with id '%d' is not found", playerID, deckID)
	}
	deck := decks[deckIdx]

	return deck, nil
}
