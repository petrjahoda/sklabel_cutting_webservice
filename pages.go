package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"time"
)

type TimePageData struct {
	Time string
}

func userBreak(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	logInfo("User Break", "Page loading...")
	var data TimePageData
	data.Time = time.Now().Format("15:04:05")
	tmpl := template.Must(template.ParseFiles("./html/user_break.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("User Break", "Page loaded")
}

func userChange(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	logInfo("Login", "Page loading...")
	var data TimePageData
	data.Time = time.Now().Format("15:04:05")
	tmpl := template.Must(template.ParseFiles("./html/user_change.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("Login", "Page loaded")
}

func idleSelect(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("Idle select", "Page loading...")
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	var data TimePageData
	data.Time = time.Now().Format("15:04:05")
	tmpl := template.Must(template.ParseFiles("./html/idle_select.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("Idle select", "Page loaded")
}

func idleRunning(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("Idle running", "Page loading...")
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	var data TimePageData
	data.Time = time.Now().Format("15:04:05")
	tmpl := template.Must(template.ParseFiles("./html/idle_running.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("Idle running", "Page loaded")
}

func cuttingEnd(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	logInfo("Cutting End", "Page loading...")
	var data TimePageData
	data.Time = time.Now().Format("15:04:05")
	tmpl := template.Must(template.ParseFiles("./html/cutting_end.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("Cutting End", "Page loaded")
}

func home(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("Home", "Home called")
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	var data TimePageData
	data.Time = time.Now().Format("15:04:05")
	tmpl := template.Must(template.ParseFiles("./html/home.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("Home", "Home ended")
}
