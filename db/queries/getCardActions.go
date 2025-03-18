package queries

import (
	"database/sql"

	"github.com/jental/freetesl-server/db"
	"github.com/jental/freetesl-server/db/models"
)

func GetCardActions() ([]*models.CardAction, error) {
	db, err := db.OpenAndTestConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT ca.card_id, ca.action_id, ca.interceptor_point_id, ca.actions_parameters_values
		FROM card_actions as ca
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cardActions := make([]*models.CardAction, 0)

	for rows.Next() {
		var cardID int
		var actionID string
		var interceptorPointID string
		var actionParametersValues sql.NullString
		if err := rows.Scan(&cardID, &actionID, &interceptorPointID, &actionParametersValues); err != nil {
			return nil, err
		}

		ca := models.CardAction{
			CardID:             cardID,
			ActionID:           actionID,
			InterceptorPointID: interceptorPointID,
		}
		if actionParametersValues.Valid {
			ca.ActionParametersValues = &actionParametersValues.String
		} else {
			ca.ActionParametersValues = nil
		}
		cardActions = append(cardActions, &ca)
	}

	return cardActions, nil
}
