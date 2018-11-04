package main

import (
	"fmt"
	"strings"
	d "google.golang.org/api/drive/v3"
	s "google.golang.org/api/sheets/v4"
)

const (
	neededSheetsId = `1RhUEggwqMzkc4HDCwQUUDRTWgFLwqMLxxCA-BERQEJM`
	readRange = `November 2018!A26:F`
	started_dir = `1GOiGA7lwqDxjK_qNEMOY8gsqbAb5TvS8`
)

func errH(err error) {
	if err != nil {
		panic(err)
	}
}

func copyTemplate(row []interface{}, drive *d.Service, sheets *s.Service) error {
	template_id := ""
	race := strings.ToLower(row[3].(string))
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
	nameSections := strings.Split(strings.TrimSpace(row[2].(string)), " ")
	fmt.Printf("%s (%s) - %s Template\n", row[0], nameSections[0], row[3])
	patch := &d.File{
		Parents: []string{started_dir},
		Name: fmt.Sprintf("%s (%s)", row[0], nameSections[0]),
	}
	new_file, err := drive.Files.Copy(template_id, patch).Do()
	if err != nil {
		return err
	}
	name := &s.ValueRange {
		MajorDimension: "ROWS",
		Range: "Character!B2:D2",
		Values: [][]interface{}{[]interface{}{row[0]}},
	}
	character := &s.ValueRange {
		MajorDimension: "ROWS",
		Range: "Character!F2:H2",
		Values: [][]interface{}{[]interface{}{row[2]}},
	}
	email := &s.ValueRange {
		MajorDimension: "ROWS",
		Range: "Character!B3:D3",
		Values: [][]interface{}{[]interface{}{row[1]}},
	}
	culture := &s.ValueRange {
		MajorDimension: "ROWS",
		Range: "Character!B4:D4",
		Values: [][]interface{}{[]interface{}{row[4]}},
	}
	native_lore := &s.ValueRange {
		MajorDimension: "ROWS",
		Range: "Character!E18",
		Values: [][]interface{}{[]interface{}{fmt.Sprintf("Native Lore: %s", row[4])}},
	}
	update := &s.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
		Data: []*s.ValueRange{name, character, email, culture, native_lore},
	}
	_, err = sheets.Spreadsheets.Values.BatchUpdate(new_file.Id, update).Do()
	return err
}

func main() {
	drive, sheets, err := InitDriveAPIs()
	resp, err := sheets.Spreadsheets.Values.Get(neededSheetsId, readRange).Do()
	errH(err)
	for _, row := range resp.Values {
		err = copyTemplate(row, drive, sheets)
		errH(err)
	}
	//if err != nil {
	//	panic(err)
	//}
	//half_fae_template := "1Y5LepZqC0W6XzdHOwzHgr1z6YZohoBfR_5BNQUzS98Q"
	//patch := &drive.File{
	//}
	//file, err := drive_srv.Files.Copy(half_fae_template, patch).Do()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(file.Id)

	//sheets_srv, err := sheets.New(client)
	//if err != nil {
	//	log.Fatalf("Unable to retrieve Sheets client: %v", err)
	//}

	//// Prints the names and majors of students in a sample spreadsheet:
	//// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	//spreadsheetId := ""
	//readRange := "November 2018!A2:E"
	//resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	//if err != nil {
	//	log.Fatalf("Unable to retrieve data from sheet: %v", err)
	//}

	//if len(resp.Values) == 0 {
	//	fmt.Println("No data found.")
	//} else {
	//	fmt.Println("Name, Major:")
	//	for _, row := range resp.Values {
	//		// Print columns A and E, which correspond to indices 0 and 4.
	//		fmt.Printf("%s, %s\n", row[0], row[4])
	//	}
	//}
}
