package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func idleSelect(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Idle select", "Page loading...")
	http.ServeFile(writer, request, "./html/idle_select.html")
	logInfo("Idle select", "Page loaded")
}
