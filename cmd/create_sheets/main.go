package main

import (
	"fmt"
	d "google.golang.org/api/drive/v3"
	s "google.golang.org/api/sheets/v4"
	"strings"
	"tm-helper/pkg/helper"
)

const (
	neededSheetsId = `1RhUEggwqMzkc4HDCwQUUDRTWgFLwqMLxxCA-BERQEJM`
	started_dir    = `1GOiGA7lwqDxjK_qNEMOY8gsqbAb5TvS8`
)

func errH(err error) {
	if err != nil {
		panic(err)
	}
}

func copyTemplate(row []string, drive *d.Service, sheets *s.Service) error {
	template_id := ""
	race := strings.ToLower(row[3])
	switch race {
	case "human":
		template_id = `1LsIuIb1__9F7ic5a4gRbO3o-Kxb3KuWt897hJjTcPXc`
	case "effendal":
		template_id = `1yzrrmNvVWdnEBkdu_BOST0apeIcMZYzW_Qz6Clo0bh8`
	case "half-celestial":
		template_id = `1t3B_xSQiHd2XGf8m0QauNONMbXu5GR3y_unT3t-w1Bg`
	case "half-demon":
		template_id = `1xxmCdm1qUrsvDbem2E8KIJJ5rVAKlGX5XMrVKCyMSN4`
	case "half-dragon":
		template_id = `1RpCKisjkQNsov48BkOzZcbkl7_ezgkC3Xy8ZPE7bYUQ`
	case "half-fae":
		template_id = `1Y5LepZqC0W6XzdHOwzHgr1z6YZohoBfR_5BNQUzS98Q`
	}
	nameSections := strings.Split(strings.TrimSpace(row[2]), " ")
	fmt.Printf("%s (%s) - %s Template\n", row[0], nameSections[0], row[3])
	patch := &d.File{
		Parents: []string{started_dir},
		Name:    fmt.Sprintf("%s (%s)", row[0], nameSections[0]),
	}
	new_file, err := drive.Files.Copy(template_id, patch).Do()
	if err != nil {
		return err
	}
	name := &s.ValueRange{
		MajorDimension: "ROWS",
		Range:          "Character!B2:D2",
		Values:         [][]interface{}{[]interface{}{row[0]}},
	}
	character := &s.ValueRange{
		MajorDimension: "ROWS",
		Range:          "Character!F2:H2",
		Values:         [][]interface{}{[]interface{}{row[2]}},
	}
	email := &s.ValueRange{
		MajorDimension: "ROWS",
		Range:          "Character!B3:D3",
		Values:         [][]interface{}{[]interface{}{row[1]}},
	}
	culture := &s.ValueRange{
		MajorDimension: "ROWS",
		Range:          "Character!B4:D4",
		Values:         [][]interface{}{[]interface{}{row[4]}},
	}
	native_lore := &s.ValueRange{
		MajorDimension: "ROWS",
		Range:          "Character!E18",
		Values:         [][]interface{}{[]interface{}{fmt.Sprintf("Native Lore: %s", row[4])}},
	}
	update := &s.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
		Data:             []*s.ValueRange{name, character, email, culture, native_lore},
	}
	_, err = sheets.Spreadsheets.Values.BatchUpdate(new_file.Id, update).Do()
	return err
}

func runNewPlayerSheets(drive *d.Service, sheets *s.Service) {
	newplayers_sheet, err := sheets.Spreadsheets.Get(neededSheetsId).IncludeGridData(true).Do()
	errH(err)
	for _, sheet := range newplayers_sheet.Sheets {
		for _, data := range sheet.Data {
			for _, row := range data.RowData {
				rowValues := []string{}
				for _, value := range row.Values {
					if value.UserEnteredValue != nil {
						rowValues = append(rowValues, value.UserEnteredValue.StringValue)
					} else {
						rowValues = append(rowValues, "")
					}
				}
				userformat := row.Values[0].UserEnteredFormat
				if len(rowValues) < 6 || len(row.Values) < 6 || strings.TrimSpace(rowValues[0]) == "" {
					continue
				}
				if (userformat == nil || userformat.BackgroundColor == nil) && rowValues[0] != "Player Name" {
					err = copyTemplate(rowValues, drive, sheets)
					errH(err)
				}
			}
		}
	}
}

func main() {
	drive, sheets, err := helper.InitDriveAPIs()
	errH(err)
	runNewPlayerSheets(drive, sheets)
	//skills, err := ReadSkills("tmskills.csv")
	//errH(err)
	//indexed, err := IndexSkills(skills)
	//errH(err)
	//err = TestSearchFunctions(indexed, sheets)
	//errH(err)
}
