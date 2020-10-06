package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"time"
)

func idleRunning(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Idle running", "Page loading...")
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	var data IdlePageData
	data.Time = time.Now().Format("15:04:05")
	tmpl := template.Must(template.ParseFiles("./html/idle_running.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("Idle running", "Page loaded")
}
