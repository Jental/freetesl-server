package queries

import (
	"fmt"

	"github.com/jental/freetesl-server/db"
)

func VerifyUser(login string, passowrdSha512 string) (bool, *int) {
	db, err := db.OpenAndTestConnection()
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
