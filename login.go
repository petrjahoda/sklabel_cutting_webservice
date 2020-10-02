package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func login(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Login", "Page loading...")
	http.ServeFile(writer, request, "./html/login.html")
	logInfo("Login", "Page loaded")
}

