package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
	di "tm-helper/pkg/docs_interface"
	"tm-helper/pkg/helper"
)

type CharacterInfo struct {
	CharacterRow           int    `sheet:"column_name:row"`
	PlayerName             string `sheet:"column_name:Player||Name"`
	PlayerEmail            string `sheet:"column_name:Player||Email"`
	PlayerEmergencyContact string `sheet:"column_name:For||safety"`
	CharacterName          string `sheet:"column_name:Character||Name"`
	CharacterRace          string `sheet:"column_name:Race"`
	CharacterNation        string `sheet:"column_name:Nation"`
	CharacterSkills        string `sheet:"column_name:Skills"`
	CharacterHistory       string `sheet:"column_name:Character||History"`
	SheetCreated           string `sheet:"column_name:Sheet||Created"`
	HistoryApproved        string `sheet:"column_name:History||Approved"`
	SheetShared            string `sheet:"column_name:Sheet||Shared"`
}

type Runner struct {
	Drive       *drive.Service
	Sheets      *sheets.Service
	Spreadsheet *di.Spreadsheet
	reload      chan bool
}

func NewRunner() (*Runner, error) {
	drive, sheets, err := helper.InitDriveAPIs()
	if err != nil {
		return nil, err
	}
	spreadsheet, err := di.LoadSpreadsheet("1avaQ8QG8bAWdbsYoKfZEnF2DPh1Wx82dJbHWwR2JiX8", sheets)
	if err != nil {
		return nil, err
	}
	toret := &Runner{
		Drive:       drive,
		Sheets:      sheets,
		Spreadsheet: spreadsheet,
		reload:      make(chan bool, 1),
	}
	go toret.Launch()
	return toret, nil
}

func (r *Runner) Launch() {
	go r.reloadLoop()
	go r.reloadTimer()
}

func (r *Runner) reloadLoop() {
	for range r.reload {
		r.Spreadsheet.Reload()
	}
}

func (r *Runner) reloadTimer() {
	t := time.Tick(5 * time.Minute)
	for range t {
		r.Reload()
	}
}

func (r *Runner) Reload() {
	r.reload <- true
}

func (r *Runner) addColumns(sheetname string, columns []string) error {
	sheet, err := r.Spreadsheet.LoadSheet(sheetname)
	if err != nil {
		return err
	}
	for _, col := range columns {
		sheet.AddColumnUpdate(col)
	}
	return sheet.ExecuteUpdates()
}

func (r *Runner) ReloadHandler(w http.ResponseWriter, req *http.Request) {
	r.Reload()
}

func (r *Runner) ViewSheetsHandler(w http.ResponseWriter, req *http.Request) {
	sheets := r.Spreadsheet.ListSheets()
	out, err := json.Marshal(sheets)
	if err != nil {
		panic(err)
	}
	w.Write(out)
}

func (r *Runner) ViewSheetHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fmt.Println("Loading Sheet:", vars["sheetname"])
	sheetname := vars["sheetname"]
	extra_cols := []string{
		"Sheet Created", "History Approved", "Sheet Shared",
	}
	r.addColumns(sheetname, extra_cols)
	sheet, err := r.Spreadsheet.LoadSheet(sheetname)
	if err != nil {
		panic(err)
	}
	characters := []CharacterInfo{}
	rc := sheet.RowCount()
	fetch := CharacterInfo{}
	for i := 2; i <= rc; i += 1 {
		err = sheet.FetchRow(i, &fetch)
		if err != nil {
			panic(err)
		}
		characters = append(characters, fetch)
	}
	out, err := json.Marshal(characters)
	if err != nil {
		panic(err)
	}
	w.Write(out)
}

func (r *Runner) ViewRowHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	sheetname := vars["sheetname"]
	row, err := strconv.Atoi(vars["row_number"])
	if err != nil {
		panic(err)
	}
	sheet, err := r.Spreadsheet.LoadSheet(sheetname)
	if err != nil {
		panic(err)
	}
	fetch := CharacterInfo{}
	err = sheet.FetchRow(row, &fetch)
	if err != nil {
		panic(err)
	}
	out, err := json.Marshal(fetch)
	if err != nil {
		panic(err)
	}
	w.Write(out)

}
