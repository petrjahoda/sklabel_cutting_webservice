package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func orderScanning(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Order Scanning", "Page loading...")
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	logInfo("Order Scanning", "Page loaded")
}
