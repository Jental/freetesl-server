package services

import (
	"maps"
	"slices"
	"sync"

	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/db/queries"
)

var allCardsCache map[int]*dbModels.Card = nil
var allCardsByNameCache map[string]*dbModels.Card = nil
var allCardCacheMtx sync.Mutex

func GetAllCards() ([]*dbModels.Card, error) {
	cardsMap, err := GetAllCardsMap()
	if err != nil {
		return nil, err
	}
	return slices.Collect(maps.Values(cardsMap)), nil
}

func GetAllCardsMap() (map[int]*dbModels.Card, error) {
	allCardCacheMtx.Lock()
	defer allCardCacheMtx.Unlock()

	if allCardsCache == nil {
		cardsFromDB, err := queries.GetAllCards()
		if err != nil {
			return nil, err
		}
		allCardsCache = make(map[int]*dbModels.Card)
		allCardsByNameCache = make(map[string]*dbModels.Card)
		for _, card := range cardsFromDB {
			allCardsCache[card.ID] = card
			allCardsByNameCache[card.Name] = card
		}
	}

	return allCardsCache, nil
}

func GetCardByName(cardName string) (*dbModels.Card, bool, error) {
	_, err := GetAllCardsMap() // to init cache
	if err != nil {
		return nil, false, err
	}

	card, exists := allCardsByNameCache[cardName]
	return card, exists, nil
}
