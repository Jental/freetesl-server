package services

import (
	"fmt"
	"log"

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
