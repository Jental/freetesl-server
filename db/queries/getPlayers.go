package queries

import (
	"github.com/jental/freetesl-server/db"
	"github.com/jental/freetesl-server/db/models"
	"github.com/jmoiron/sqlx"
)

func GetPlayers() ([]*models.Player, error) {
	db, err := db.OpenAndTestConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var query string
	var args []interface{}

	query = `
		SELECT p.id, p.display_name, p.avatar_name
		FROM players as p
		ORDER BY p.id ASC
	`
	args = make([]interface{}, 0)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.Player, 0)

	for rows.Next() {
		var id int
		var displayName string
		var avatarName string
		if err := rows.Scan(&id, &displayName, &avatarName); err != nil {
			return nil, err
		}

		result = append(result, &models.Player{
			ID:          id,
			DisplayName: displayName,
			AvatarName:  avatarName,
		})
	}

	return result, nil
}

func GetPlayersByIDs(playerIDs []int) (map[int]*models.Player, error) {
	db, err := db.OpenAndTestConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var query string
	var args []interface{}

	query, args, err = sqlx.In(`
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
