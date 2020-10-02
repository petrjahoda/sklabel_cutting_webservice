package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func orderError(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Order Error", "Page loading...")
	http.ServeFile(writer, request, "./html/order_error.html")
	logInfo("Order Error", "Page loaded")
}