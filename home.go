package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func home(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Home", "Page loading...")
	http.ServeFile(writer, request, "./html/home.html")
	logInfo("Home", "Page loaded")
}
