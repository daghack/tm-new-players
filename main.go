package main

import (
	"fmt"
	"context"
	"io/ioutil"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func errH(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	data, err := ioutil.ReadFile("sheets_creds.json")
	errH(err)
	conf, err := google.JWTConfigFromJSON(data, sheets.SpreadsheetsScope)
	errH(err)
	client := conf.Client(context.TODO())
	srv, err := sheets.New(client)
	errH(err)

	spreadsheetId := "1RhUEggwqMzkc4HDCwQUUDRTWgFLwqMLxxCA-BERQEJM"
	spreadsheetId = "1GBhYD0GxraFtJc2EY3noQwSUA0o8w4YfhxrNabS9bFE"
	readRange := "Class Data!A2:E"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	errH(err)
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for _, row := range resp.Values {
			fmt.Println(row)
		}
	}
}
