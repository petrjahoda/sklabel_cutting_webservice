package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func idleRunning(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Idle running", "Page loading...")
	http.ServeFile(writer, request, "./html/idle_running.html")
	logInfo("Idle running", "Page loaded")
}
