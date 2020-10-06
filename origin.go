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

func origin(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ipAddress := strings.Split(request.Host, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Checking initial conditions")
	terminalId, workplaceCode, ipAddressHasAssignedTerminal := checkIpAddress(deviceName, ipAddress[0])
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	if !ipAddressHasAssignedTerminal {
		var data OriginPage
		data.Information = "Ip adresa " + ipAddress[0] + " nemá přiřazeno žádné pracoviště."
		tmpl := template.Must(template.ParseFiles("./html/origin.html"))
		_ = tmpl.Execute(writer, data)
		logInfo(deviceName, "Checking ended, ip address not assigned")
		return
	}
	userId, user, userIsLogged := checkIfUserIsLoggedForTerminalId(deviceName, terminalId)
	if !userIsLogged {
		var data OriginPage
		data.Information = "Na pracovišti není přihlášen žádný operátor"
		tmpl := template.Must(template.ParseFiles("./html/origin.html"))
		_ = tmpl.Execute(writer, data)
		logInfo(deviceName, "Checking ended, user not assigned")
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
	logInfo(deviceName, "Checking ended, everything OK")
	return

}
