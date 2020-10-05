package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"time"
)

type IdlePageData struct {
	Time string
}

func idleSelect(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Idle select", "Page loading...")
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	var data IdlePageData
	data.Time = time.Now().Format("15:04:05")
	tmpl := template.Must(template.ParseFiles("./html/idle_select.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("Idle select", "Page loaded")
}
