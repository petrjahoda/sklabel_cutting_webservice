package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"time"
)

func cuttingEnd(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	logInfo("Cutting End", "Page loading...")
	var data HomePageData
	data.Time = time.Now().Format("15:04:05")
	tmpl := template.Must(template.ParseFiles("./html/cutting_end.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("Cutting End", "Page loaded")
}
