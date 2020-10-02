package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"time"
)

type HomePageData struct {
	Time string
}

func home(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Home", "Home called")
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	var data HomePageData
	data.Time = time.Now().Format("15:04:05")
	tmpl := template.Must(template.ParseFiles("./html/home.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("Home", "Home ended")
}
