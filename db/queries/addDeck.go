package queries

import (
	"fmt"
	"strings"

	"github.com/jental/freetesl-server/db"
	"github.com/jental/freetesl-server/db/models"
)

func AddDeck(request *models.AddDeckDbRequest) (int, error) {
	db, err := db.OpenAndTestConnection()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	rows, err := db.Query(`
		INSERT INTO decks (name, player_id, avatar_name)
		VALUES ($1, $2, $3)
		RETURNING id
	`, request.Name, request.PlayerID, request.AvatarName)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	if !rows.Next() {
		return -1, fmt.Errorf("AddDeck: id is expected to be returned")
	}

	var id int
	err = rows.Scan(&id)
	if err != nil {
		return -1, err
	}

	valuesStrings := make([]string, 0)
	parameters := []any{id}
	parameterNum := 2
	for cardID, count := range request.Cards {
		valuesStrings = append(valuesStrings, fmt.Sprintf("($1, $%d, $%d)", parameterNum, parameterNum+1))
		parameters = append(parameters, cardID, count)
		parameterNum = parameterNum + 2
	}

	valuesString := strings.Join(valuesStrings, ", ")
	query := fmt.Sprintf(`
		INSERT INTO deck_cards (deck_id, card_id, count)
		VALUES %s
	`, valuesString)
	_, err = db.Query(query, parameters...)
	if err != nil {
		return id, err
	}

	return id, nil
}
