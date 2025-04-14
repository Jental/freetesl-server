package queries

import (
	"context"
	"fmt"

	"github.com/jental/freetesl-server/db"
)

func DeleteDeck(ctx *context.Context, id int, playerID int) error {
	db, err := db.OpenAndTestConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT count(*) FROM decks WHERE id = $1 and player_id = $2`, id, playerID)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf("DeleteDeck: count is expected to be returned")
	}
	var count int
	err = rows.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("DeleteDeck: no deck found")
	}

	tx, err := db.BeginTx(*ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = db.Exec(`DELETE FROM deck_cards WHERE deck_id = $1`, id)
	if err != nil {
		return err
	}

	_, err = db.Exec(`DELETE FROM decks WHERE id = $1 and player_id = $2`, id, playerID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
