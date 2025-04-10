package queries

import (
	"github.com/jental/freetesl-server/db"
	"github.com/jental/freetesl-server/db/models"
	"golang.org/x/exp/maps"
)

func GetClasses() ([]*models.CardClass, error) {
	db, err := db.OpenAndTestConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT c.id as id, c.name as name, a.id, a.name
		FROM classes as c
		INNER JOIN classes_to_attributes as c2a ON c2a.class_id = c.id
		INNER JOIN attributes as a ON a.id = c2a.attribute_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	classes := make(map[byte]*models.CardClass)

	for rows.Next() {
		var id byte
		var name string
		var attributeID int
		var attributeName string
		if err := rows.Scan(&id, &name, &attributeID, &attributeName); err != nil {
			return nil, err
		}

		cardClass, exists := classes[id]
		if !exists {
			var newClass = models.CardClass{ID: id, Name: name, Attributes: make([]*models.Attribute, 0)}
			cardClass = &newClass
			classes[id] = &newClass
		}

		attribute := models.Attribute{
			ID:   attributeID,
			Name: attributeName,
		}
		cardClass.Attributes = append(cardClass.Attributes, &attribute)
	}

	return maps.Values(classes), nil
}
