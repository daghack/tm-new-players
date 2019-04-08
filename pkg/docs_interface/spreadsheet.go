package docs_interface

import (
	"fmt"
	s "google.golang.org/api/sheets/v4"
)

type Spreadsheet struct {
	service     *s.Service
	spreadsheet *s.Spreadsheet
}

func LoadSpreadsheet(driveId string, service *s.Service) (*Spreadsheet, error) {
	spreadsheet, err := service.Spreadsheets.Get(driveId).IncludeGridData(true).Do()
	if err != nil {
		return nil, err
	}
	return &Spreadsheet{
		service:     service,
		spreadsheet: spreadsheet,
	}, nil
}

func (this *Spreadsheet) Id() string {
	return this.spreadsheet.SpreadsheetId
}

func (this *Spreadsheet) LoadSheet(sheetName string) (*Sheet, error) {
	for _, sheet := range this.spreadsheet.Sheets {
		if sheet.Properties.Title == sheetName {
			return &Sheet{
				parent:  this,
				service: this.service,
				sheet:   sheet,
				name:    sheetName,
			}, nil
		}
	}
	return nil, fmt.Errorf("No sheet named '%s' in spreadsheet.", sheetName)
}

func (this *Spreadsheet) Reload() error {
	spreadsheet, err := this.service.Spreadsheets.Get(this.Id()).IncludeGridData(true).Do()
	if err != nil {
		return err
	}
	this.spreadsheet = spreadsheet
	return nil
}

func (this *Spreadsheet) ListSheets() []string {
	toret := []string{}
	for _, sheet := range this.spreadsheet.Sheets {
		toret = append(toret, sheet.Properties.Title)
	}
	return toret
}
