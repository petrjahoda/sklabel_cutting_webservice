package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type TimePageData struct {
	Time string
}

func userBreak(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ipAddress := strings.Split(request.RemoteAddr, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Sending user break page")
	var data TimePageData
	data.Time = time.Now().Format("15:04")
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	tmpl := template.Must(template.ParseFiles("./html/user_break.html"))
	_ = tmpl.Execute(writer, data)
	logInfo(deviceName, "User break page sent")
}

func userChange(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ipAddress := strings.Split(request.RemoteAddr, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Sending user change page")
	var data TimePageData
	data.Time = time.Now().Format("15:04")
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	tmpl := template.Must(template.ParseFiles("./html/user_change.html"))
	_ = tmpl.Execute(writer, data)
	logInfo(deviceName, "User change page sent")
}

func idleSelect(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ipAddress := strings.Split(request.RemoteAddr, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Sending idle select page")
	var data TimePageData
	data.Time = time.Now().Format("15:04")
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	tmpl := template.Must(template.ParseFiles("./html/idle_select.html"))
	_ = tmpl.Execute(writer, data)
	logInfo(deviceName, "DataIdle select page sent")

}

func idleRunning(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ipAddress := strings.Split(request.RemoteAddr, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Sending idle running page")
	var data TimePageData
	data.Time = time.Now().Format("15:04")
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	tmpl := template.Must(template.ParseFiles("./html/idle_running.html"))
	_ = tmpl.Execute(writer, data)
	logInfo(deviceName, "DataIdle running page sent")

}

func cuttingEnd(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ipAddress := strings.Split(request.RemoteAddr, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Sending cutting end page")
	var data TimePageData
	data.Time = time.Now().Format("15:04")
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	tmpl := template.Must(template.ParseFiles("./html/cutting_end.html"))
	_ = tmpl.Execute(writer, data)
	logInfo(deviceName, "Cutting end page sent")
}

func home(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ipAddress := strings.Split(request.RemoteAddr, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Sending home page")
	var data TimePageData
	data.Time = time.Now().Format("15:04")
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	tmpl := template.Must(template.ParseFiles("./html/home.html"))
	_ = tmpl.Execute(writer, data)
	logInfo(deviceName, "Home page sent")
}
