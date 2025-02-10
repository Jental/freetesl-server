package main

import (
	"database/sql"

	"github.com/jental/freetesl-server/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/exp/maps"
)

const CONNECTION_STRING = "host=localhost port=5432 user=postgres password=]Hy)*58)Np-2LrC9hD-( dbname=tesl sslmode=disable"

func getDecks(playerID int) ([]models.Deck, error) {
	db, err := sql.Open("postgres", CONNECTION_STRING)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
		SELECT d.id as id, d.name as name, c.id as card_id, c.name as card_name, c.power as card_power, c.defence as card_defence, c.cost as card_cost, dc.count as count
		FROM decks as d 
		INNER JOIN deck_cards as dc ON dc.deck_id = d.id
		INNER JOIN cards as c ON c.id = dc.card_id
		WHERE d.player_id = $1
	`, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	decks := make(map[int]models.Deck)

	for rows.Next() {
		var id int
		var name string
		var cardID int
		var cardName string
		var cardPower int
		var cardDefence int
		var cardCost int
		var count int
		if err := rows.Scan(&id, &name, &cardID, &cardName, &cardPower, &cardDefence, &cardCost, &count); err != nil {
			return nil, err
		}

		deck, exists := decks[id]
		if !exists {
			deck = models.Deck{ID: id, Name: name, Cards: make(map[int]models.CardWithCount)}
			decks[id] = deck
		}

		deck.Cards[cardID] = models.CardWithCount{
			Card: models.Card{
				ID:          cardID,
				Name:        cardName,
				Description: "dscr",
				Power:       cardPower,
				Health:      cardDefence,
				Cost:        cardCost,
			},
			Count: count,
		}
	}

	return maps.Values(decks), nil
}

func getPlayers(playerIDs []int) (map[int]*models.Player, error) {
	db, err := sql.Open("postgres", CONNECTION_STRING)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// var playerIDsStr = strings.Replace(strings.Trim(fmt.Sprint(playerIDs), "[]"), " ", ", ", -1)
	query, args, err := sqlx.In(`
		SELECT p.id, p.display_name, p.avatar_name
		FROM players as p 
		WHERE p.id in (?)
	`, playerIDs)
	if err != nil {
		return nil, err
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)
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
