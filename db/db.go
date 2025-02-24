package db

import (
	"database/sql"
	"fmt"

	"github.com/jental/freetesl-server/db/enums"
	"github.com/jental/freetesl-server/db/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
)

const CONNECTION_STRING = "host=localhost port=5432 user=postgres password=]Hy)*58)Np-2LrC9hD-( dbname=tesl sslmode=disable"

func openAndTestConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", CONNECTION_STRING)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		defer db.Close()
		return nil, err
	}

	return db, nil
}

func GetAllCards() ([]*models.Card, error) {
	db, err := openAndTestConnection()
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

func GetDecks(playerID int) ([]*models.Deck, error) {
	db, err := openAndTestConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT d.id as id, d.name as name, c.id as card_id, c.name as card_name, c.power as card_power, c.health as card_health, c.cost as card_cost, dc.count as count
		FROM decks as d 
		INNER JOIN deck_cards as dc ON dc.deck_id = d.id
		INNER JOIN cards as c ON c.id = dc.card_id
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
		var cardID int
		var cardName string
		var cardPower int
		var cardHealth int
		var cardCost int
		var count int
		if err := rows.Scan(&id, &name, &cardID, &cardName, &cardPower, &cardHealth, &cardCost, &count); err != nil {
			return nil, err
		}

		deck, exists := decks[id]
		if !exists {
			var newDeck = models.Deck{ID: id, Name: name, Cards: make(map[int]models.CardWithCount)}
			deck = &newDeck
			decks[id] = &newDeck
		}

		deck.Cards[cardID] = models.CardWithCount{
			Card: models.Card{
				ID:     cardID,
				Name:   cardName,
				Power:  cardPower,
				Health: cardHealth,
				Cost:   cardCost,
			},
			Count: count,
		}
	}

	return maps.Values(decks), nil
}

func GetPlayers(playerIDs []int) (map[int]*models.Player, error) {
	db, err := openAndTestConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var query string
	var args []interface{}
	if playerIDs != nil {
		query, args, err = sqlx.In(`
			SELECT p.id, p.display_name, p.avatar_name
			FROM players as p 
			WHERE p.id in (?)
			ORDER BY p.id ASC
		`, playerIDs)
		if err != nil {
			return nil, err
		}
		query = sqlx.Rebind(sqlx.DOLLAR, query)

	} else {
		query = `
			SELECT p.id, p.display_name, p.avatar_name
			FROM players as p
			ORDER BY p.id ASC
		`
		args = make([]interface{}, 0)
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]*models.Player)

	for rows.Next() {
		var id int
		var displayName string
		var avatarName string
		if err := rows.Scan(&id, &displayName, &avatarName); err != nil {
			return nil, err
		}

		result[id] = &models.Player{
			ID:          id,
			DisplayName: displayName,
			AvatarName:  avatarName,
		}
	}

	return result, nil
}

func VerifyUser(login string, passowrdSha512 string) (bool, *int) {
	db, err := openAndTestConnection()
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT id
		FROM players
		WHERE login = $1 AND password = $2
	`, login, passowrdSha512)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	defer rows.Close()

	exists := rows.Next()
	if !exists {
		return false, nil
	}

	var userID int
	err = rows.Scan(&userID)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}

	return true, &userID
}
