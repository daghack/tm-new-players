package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	runner, err := NewRunner()
	if err != nil {
		panic(err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/reload", runner.ReloadHandler)
	router.HandleFunc("/sheets", runner.ViewSheetsHandler)
	router.HandleFunc("/sheet/{sheetname}", runner.ViewSheetHandler)
	router.HandleFunc("/sheet/{sheetname}/{row_number}", runner.ViewRowHandler)
	panic(http.ListenAndServe(":8080", router))
}
