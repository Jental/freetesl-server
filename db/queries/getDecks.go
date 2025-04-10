package queries

import (
	"github.com/jental/freetesl-server/db"
	"github.com/jental/freetesl-server/db/models"
	"golang.org/x/exp/maps"
)

func GetDecks(playerID int) ([]*models.Deck, error) {
	db, err := db.OpenAndTestConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT d.id as id, d.name as name, d.avatar_name, dc.card_id as card_id, dc.count as count
		FROM decks as d 
		INNER JOIN deck_cards as dc ON dc.deck_id = d.id
		WHERE d.player_id = $1
	`, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	decks := make(map[int]*models.Deck)

	for rows.Next() {
		var id int
		var name string
		var avatarName string
		var cardID int
		var count int
		if err := rows.Scan(&id, &name, &avatarName, &cardID, &count); err != nil {
			return nil, err
		}

		deck, exists := decks[id]
		if !exists {
			var newDeck = models.Deck{ID: id, Name: name, AvatarName: avatarName, Cards: make(map[int]int)}
			deck = &newDeck
			decks[id] = &newDeck
		}

		deck.Cards[cardID] = count
	}

	return maps.Values(decks), nil
}
