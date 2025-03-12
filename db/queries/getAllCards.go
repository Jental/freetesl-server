package queries

import (
	"database/sql"

	"github.com/jental/freetesl-server/db"
	"github.com/jental/freetesl-server/db/enums"
	"github.com/jental/freetesl-server/db/models"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
)

func GetAllCards() ([]*models.Card, error) {
	db, err := db.OpenAndTestConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT c.id as id, c.name, c.power, c.health, c.cost, c.type_id, c.class_id, cr.race_id, ck.keyword_id
		FROM cards as c
		LEFT JOIN card_races as cr ON cr.card_id = c.id
		LEFT JOIN card_keywords as ck ON ck.card_id = c.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards = make(map[int]*models.Card)
	var keywords = make(map[int]map[byte]bool) // sort of hashset
	var races = make(map[int]map[byte]bool)

	for rows.Next() {
		var id int
		var name string
		var power int
		var health int
		var cost int
		var typeID int
		var classID int
		var raceID sql.NullByte
		var keywordID sql.NullByte
		if err := rows.Scan(&id, &name, &power, &health, &cost, &typeID, &classID, &raceID, &keywordID); err != nil {
			return nil, err
		}

		_, exists := cards[id]
		if !exists {
			var newCard = models.Card{
				ID:      id,
				Name:    name,
				Power:   power,
				Health:  health,
				Cost:    cost,
				Type:    enums.CardType(typeID),
				ClassID: byte(classID),
			}
			cards[id] = &newCard
		}

		cardRaces, exists := races[id]
		if !exists {
			cardRaces = make(map[byte]bool)
			races[id] = cardRaces
		}
		if raceID.Valid {
			cardRaces[raceID.Byte] = true
		}

		cardKeywords, exists := keywords[id]
		if !exists {
			cardKeywords = make(map[byte]bool)
			keywords[id] = cardKeywords
		}
		if keywordID.Valid {
			cardKeywords[keywordID.Byte] = true
		}
	}

	var cardValues = maps.Values(cards)

	for _, card := range cardValues {
		cardRaces := races[card.ID]
		card.Races = maps.Keys(cardRaces)

		cardKeywords := keywords[card.ID]
		card.Keywords = lo.Map(maps.Keys(cardKeywords), func(kw byte, _ int) enums.CardKeyword { return enums.CardKeyword(kw) })
	}

	return cardValues, nil
}
