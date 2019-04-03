package main

import (
	"fmt"
	di "tm-helper/pkg/docs_interface"
	"tm-helper/pkg/helper"
)

func init() {
}

func errH(err error) {
	if err != nil {
		panic(err)
	}
}

type CharacterInfo struct {
	PlayerName             string `sheet:"column_name:Player||Name"`
	PlayerEmail            string `sheet:"column_name:Player||Email"`
	PlayerEmergencyContact string `sheet:"column_name:For||safety"`
	CharacterName          string `sheet:"column_name:Character||Name"`
	CharacterRace          string `sheet:"column_name:Race"`
	CharacterNation        string `sheet:"column_name:Nation"`
	CharacterSkills        string `sheet:"column_name:Skills"`
	CharacterHistory       string `sheet:"column_name:Character||History"`
}

func main() {
	_, sheets, err := helper.InitDriveAPIs()
	errH(err)

	spreadsheet, err := di.LoadSpreadsheet("1avaQ8QG8bAWdbsYoKfZEnF2DPh1Wx82dJbHWwR2JiX8", sheets)
	errH(err)

	sheet, err := spreadsheet.LoadSheet("TestData 2019")
	errH(err)

	spreadsheet, err = di.LoadSpreadsheet("1avaQ8QG8bAWdbsYoKfZEnF2DPh1Wx82dJbHWwR2JiX8", sheets)
	sheet, err = spreadsheet.LoadSheet("TestData 2019")
	errH(err)

	fetch := &CharacterInfo{}
	err = sheet.FetchRow(160, fetch)
	errH(err)
	fmt.Printf("%+v\n", fetch)

	rc := sheet.RowCount()
	for i := 2; i <= rc; i += 1 {
		err = sheet.FetchRow(i, fetch)
		errH(err)
		fmt.Printf("%+v\n", fetch)
	}

	//	submission_sheet, err := sheets.Spreadsheets.Get("1avaQ8QG8bAWdbsYoKfZEnF2DPh1Wx82dJbHWwR2JiX8").IncludeGridData(true).Do()
	//	errH(err)
	//	for _, sheet := range submission_sheet.Sheets {
	//		fmt.Println(sheet.Properties.Title)
	//		if sheet.Properties.Title == "TestData 2019" {
	//			for _, value := range sheet.Data[0].RowData[0].Values {
	//				if value.UserEnteredValue != nil || strings.TrimSpace(value.UserEnteredValue.StringValue) != "" {
	//					fmt.Printf("%+v\n", value.UserEnteredValue)
	//				} else {
	//					fmt.Println("nil")
	//				}
	//			}
	//		}
	//	}
}
