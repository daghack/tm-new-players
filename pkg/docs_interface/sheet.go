package docs_interface

import (
	"fmt"
	s "google.golang.org/api/sheets/v4"
	"reflect"
	"strings"
)

type Sheet struct {
	parent  *Spreadsheet
	service *s.Service
	sheet   *s.Sheet
	name    string
	updates map[string]string
}

func (this *Sheet) ExecuteUpdates() error {
	update := &s.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
		Data:             []*s.ValueRange{},
	}
	for cell, value := range this.updates {
		update.Data = append(update.Data, &s.ValueRange{
			MajorDimension: "ROWS",
			Range:          cell,
			Values:         [][]interface{}{[]interface{}{value}},
		})
	}
	_, err := this.service.Spreadsheets.Values.BatchUpdate(this.parent.Id(), update).Do()
	if err != nil {
		return err
	}
	this.updates = map[string]string{}
	err = this.parent.Reload()
	if err != nil {
		return err
	}
	s, err := this.parent.LoadSheet(this.name)
	if err != nil {
		return err
	}
	this.sheet = s.sheet
	return nil
}

func (this *Sheet) AddUpdate(cell, content string) {
	if this.updates == nil {
		this.updates = map[string]string{}
	}
	this.updates[cell] = content
}

func (this *Sheet) FetchColumnIndex(columnName string) int {
	for i, value := range this.sheet.Data[0].RowData[0].Values {
		if value.UserEnteredValue != nil && strings.HasPrefix(value.UserEnteredValue.StringValue, columnName) {
			return i
		}
	}
	return -1
}

func (this *Sheet) AddRowUpdate(row int, patchInfo interface{}) error {
	valueOf := reflect.ValueOf(patchInfo)
	typeOf := valueOf.Type()
	for typeOf.Kind() == reflect.Ptr {
		valueOf = reflect.Indirect(valueOf)
		typeOf = valueOf.Type()
	}
	fieldCount := typeOf.NumField()
	for i := 0; i < fieldCount; i += 1 {
		field := typeOf.Field(i)
		if field.Type.Kind() != reflect.String {
			continue
		}
		tagSet := getFieldTags(field)
		if col, ok := tagSet["column_name"]; ok {
			index := this.FetchColumnIndex(col)
			if index < 0 {
				if _, ok := tagSet["add_column"]; ok {
					index = this.AddColumnUpdate(col)
				}
			}
			if index >= 0 {
				cellData := valueOf.FieldByName(field.Name).String()
				cell := fmt.Sprintf("%s!%s%d", this.sheet.Properties.Title, string(rune('A'+index)), row)
				this.AddUpdate(cell, cellData)
			}
		}
	}
	return nil
}

func getFieldTags(field reflect.StructField) map[string]string {
	sheetTag := field.Tag.Get("sheet")
	tags := strings.Split(sheetTag, " ")
	tagSet := map[string]string{}
	for _, tag := range tags {
		tag = strings.ReplaceAll(tag, "||", " ")
		keyval := strings.Split(tag, ":")
		if len(keyval) == 2 {
			tagSet[strings.TrimSpace(keyval[0])] = strings.TrimSpace(keyval[1])
		} else if len(keyval) == 1 {
			tagSet[strings.TrimSpace(keyval[0])] = ""
		}
	}
	return tagSet
}

func (this *Sheet) AddColumnUpdate(columnName string) int {
	column := 0
	onblank := false
	for i, value := range this.sheet.Data[0].RowData[0].Values {
		if value.UserEnteredValue == nil || strings.TrimSpace(value.UserEnteredValue.StringValue) == "" {
			if !onblank {
				column = i
				onblank = true
			}
		} else {
			if value.UserEnteredValue != nil && strings.HasPrefix(value.UserEnteredValue.StringValue, columnName) {
				return i
			}
			column = len(this.sheet.Data[0].RowData[0].Values)
			onblank = false
		}
	}
	cell := fmt.Sprintf("%s!%s1", this.sheet.Properties.Title, string(rune('A'+column)))
	for _, ok := this.updates[cell]; ok; _, ok = this.updates[cell] {
		cell = fmt.Sprintf("%s!%s1", this.sheet.Properties.Title, string(rune('A'+column)))
		column += 1
	}
	fmt.Println(cell)
	this.AddUpdate(cell, columnName)
	return column
}

func (this *Sheet) UpdateRow(row int, patch interface{}) error {
	err := this.AddRowUpdate(row, patch)
	if err != nil {
		return err
	}
	return this.ExecuteUpdates()
}

func (this *Sheet) FetchRow(row int, s interface{}) error {
	if row < 1 || row > len(this.sheet.Data[0].RowData) {
		return fmt.Errorf("Invalid Index: %d", row)
	}
	valueOf := reflect.ValueOf(s)
	typeOf := valueOf.Type()
	for typeOf.Kind() == reflect.Ptr {
		valueOf = reflect.Indirect(valueOf)
		typeOf = valueOf.Type()
	}
	fieldCount := typeOf.NumField()
	for i := 0; i < fieldCount; i += 1 {
		field := typeOf.Field(i)
		tagSet := getFieldTags(field)
		colname, ok := tagSet["column_name"]
		if colname == "row" && field.Type.Kind() == reflect.Int {
			valueOf.Field(i).SetInt(int64(row))
		}
		if !ok || field.Type.Kind() != reflect.String {
			continue
		}
		col := this.FetchColumnIndex(colname)
		if col < 0 {
			continue
		}
		cellData := ""
		if len(this.sheet.Data[0].RowData[row-1].Values) > col {
			if this.sheet.Data[0].RowData[row-1].Values[col].UserEnteredValue != nil {
				cellData = this.sheet.Data[0].RowData[row-1].Values[col].UserEnteredValue.StringValue
			}
		}
		valueOf.Field(i).SetString(cellData)
	}
	return nil
}

func (this *Sheet) RowCount() int {
	for i, row := range this.sheet.Data[0].RowData {
		if len(row.Values) == 0 {
			return i
		}
		allnil := true
		for _, v := range row.Values {
			if v.UserEnteredValue != nil {
				allnil = false
				break
			}
		}
		if allnil {
			return i
		}
	}
	return len(this.sheet.Data[0].RowData)
}
