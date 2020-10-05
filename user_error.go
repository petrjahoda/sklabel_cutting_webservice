package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func userError(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("User Error", "Page loading...")
	http.ServeFile(writer, request, "./html/user_break.html")
	logInfo("User Error", "Page loaded")
}
