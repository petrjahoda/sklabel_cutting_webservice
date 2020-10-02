package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func orderScanning(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Order Scanning", "Page loading...")
	http.ServeFile(writer, request, "./html/order_scanning.html")
	logInfo("Order Scanning", "Page loaded")
}