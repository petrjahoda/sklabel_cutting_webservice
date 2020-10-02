package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func entryPcs(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Entry Pcs", "Page loading...")
	http.ServeFile(writer, request, "./html/entry_pcs.html")
	logInfo("Entry Pcs", "Page loaded")
}
