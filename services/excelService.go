package services

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	dbModels "github.com/jental/freetesl-server/db/models"
	"github.com/jental/freetesl-server/models"
	"github.com/xuri/excelize/v2"
)

const sheetName = "Deck"

func ExportDeckToExcel(deck *models.Deck) ([]byte, error) {
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Create a new sheet.
	index, err := file.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	file.SetActiveSheet(index)
	file.DeleteSheet("Sheet1")

	file.SetCellValue(sheetName, "A1", "Name:")
	file.SetCellValue(sheetName, "B1", deck.Name)

	file.SetCellValue(sheetName, "A2", "Image:")
	file.SetCellValue(sheetName, "B2", deck.AvatarName)
	file.SetCellValue(sheetName, "D2", "Avaliable values:")
	file.SetCellValue(sheetName, "E2", "crdl_04_119_avatar_png")
	file.SetCellValue(sheetName, "F2", "DBH_NPC_CRDL_02_022_avatar_png")

	file.SetCellValue(sheetName, "A3", "Cards:")
	rowIdx := 4
	for _, card := range deck.Cards {
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIdx), card.Card.Name)
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIdx), card.Count)
		rowIdx = rowIdx + 1
	}

	file.SetColWidth(sheetName, "A", "B", 30)
	file.SetColWidth(sheetName, "D", "F", 30)

	buff, err := file.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func ImportDeckFromExcel(playerID int, content []byte) (*dbModels.AddDeckDbRequest, error) {
	reader := bytes.NewReader(content)
	file, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	deckName, err := file.GetCellValue(sheetName, "B1")
	if err != nil {
		return nil, err
	}

	avatarName, err := file.GetCellValue(sheetName, "B2")
	if err != nil {
		return nil, err
	}

	deck := dbModels.AddDeckDbRequest{
		PlayerID:   playerID,
		Name:       deckName,
		AvatarName: avatarName,
		Cards:      make(map[int]int),
	}

	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	for i, row := range rows {
		if i <= 3 {
			continue
		}
		if len(row) < 2 {
			return nil, fmt.Errorf("invalid data int a row '%d'. Expected two cells: card name and count", i)
		}

		cardName := row[0]
		card, exists, err := GetCardByName(cardName)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("card with name '%s' is not found", cardName)
		}

		countStr := row[1]
		count, err := strconv.Atoi(countStr)
		if err != nil {
			return nil, err
		}

		deck.Cards[card.ID] = count
	}

	return &deck, nil
}
