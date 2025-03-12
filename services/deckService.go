package services

import (
	"errors"
	"fmt"

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

	decks := lo.Map(decksFromDB, func(dbDeck *dbModels.Deck, _ int) *models.Deck {
		deck := models.Deck{
			ID:       dbDeck.ID,
			Name:     dbDeck.Name,
			PlayerID: dbDeck.PlayerID,
			Cards: lo.MapToSlice(dbDeck.Cards, func(cardID int, count int) *models.CardWithCount {
				card, exists := allCards[cardID]
				if !exists {
					cardWithCount := models.CardWithCount{
						Card:  nil,
						Count: -cardID, // hack to pass it to error
					}
					return &cardWithCount
				}

				cardWithCount := models.CardWithCount{
					Card:  card,
					Count: count,
				}

				return &cardWithCount
			}),
		}
		return &deck
	})

	decksWithMissingCards := make([]struct {
		deckID int
		cardID int
	}, 0)
	for _, deck := range decks {
		for _, card := range deck.Cards {
			if card.Card == nil {
				decksWithMissingCards = append(decksWithMissingCards, struct {
					deckID int
					cardID int
				}{deckID: deck.ID, cardID: -card.Count})
			}
		}
	}

	if len(decksWithMissingCards) > 0 {
		errs := lo.Map(decksWithMissingCards, func(el struct {
			deckID int
			cardID int
		}, _ int) error {
			return fmt.Errorf("deck with id '%d' have unexisting card: '%d'", el.deckID, el.cardID)
		})
		return nil, errors.Join(errs...)
	}

	return decks, nil
}
