package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type OrderScanningPage struct {
	WorkplaceCode string
	User          string
	Time          string
	DeviceId      string
	UserId        string
}
type OriginPage struct {
	Information string
}

func origin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Origin", "Page loading...")
	ipAddress := strings.Split(request.Host, ":")
	terminalId, workplaceCode, ipAddressHasAssignedTerminal := checkIpAddress(ipAddress[0])
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	if !ipAddressHasAssignedTerminal {
		var data OriginPage
		data.Information = "Ip adresa " + ipAddress[0] + " nemá přiřazeno žádné pracoviště."
		tmpl := template.Must(template.ParseFiles("./html/origin.html"))
		_ = tmpl.Execute(writer, data)
		logInfo("Origin", "Page loaded")
		return
	}
	userId, user, userIsLogged := checkIfUserIsLoggedForTerminalId(terminalId)
	if !userIsLogged {
		var data OriginPage
		data.Information = "Na pracovišti není přihlášen žádný operátor"
		tmpl := template.Must(template.ParseFiles("./html/origin.html"))
		_ = tmpl.Execute(writer, data)
		logInfo("Origin", "Page loaded")
		return
	}
	var data OrderScanningPage
	data.WorkplaceCode = workplaceCode
	data.User = user
	data.Time = time.Now().Format("15:04:05")
	data.DeviceId = strconv.Itoa(terminalId)
	data.UserId = strconv.Itoa(userId)
	tmpl := template.Must(template.ParseFiles("./html/order_scanning.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("Origin", "Page loaded")
	return

}

func checkIfUserIsLoggedForTerminalId(terminalId int) (int, string, bool) {
	//TODO: check user logging
	logInfo("Origin", "Checking if user is logged for terminal id: "+strconv.Itoa(terminalId))
	logInfo("Origin", "User is logged")
	return 23, "Petr Jahoda", true
}

func checkIpAddress(ipAddress string) (int, string, bool) {
	//TODO: check terminal assign
	logInfo("Origin", "Checking ip address: "+ipAddress)
	logInfo("Origin", "Ip address assigned to terminal")
	return 1, "E235", true
}
