package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"strings"
)

type ZapsiResponseData struct {
	Data     string
	UserId   int
	UserName string
}

type ZapsiData struct {
	Data string
}

type Idle struct {
	Id   int
	Name string
}

type Idles struct {
	Data []Idle
}

type OrderData struct {
	Order    string
	DeviceId string
	UserId   string
}

type IdleData struct {
	Order    string
	IdleId   string
	DeviceId string
	UserId   string
}

func checkIfUserIsLoggedForTerminalId(deviceName string, terminalId int) (int, string, bool) {
	//TODO: check if user is logged in terminal input user (nebo jak)
	logInfo(deviceName, "Checking if user is logged for terminal id: "+strconv.Itoa(terminalId))
	logInfo(deviceName, "User is logged")
	return 23, "Petr Jahoda", true
}

func checkIpAddress(deviceName string, ipAddress string) (int, string, bool) {
	//TODO: check assigned terminal to this ip
	logInfo(deviceName, "Checking ip address: "+ipAddress)
	logInfo(deviceName, "Ip address assigned to terminal")
	return 1, "E235", true
}

func startOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ipAddress := strings.Split(request.Host, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Save order in Zapsi called")
	var data OrderData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", err.Error())
		return
	}
	logInfo(deviceName, "Order: "+data.Order+"; userId:"+data.UserId+"; deviceId: "+data.DeviceId)
	//TODO: StartOrderInZapsi(data)
	logInfo(deviceName, "Save order in Zapsi finished, everything ok")
}

func checkUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ipAddress := strings.Split(request.Host, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Check user in Zapsi called")
	var data ZapsiData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError(deviceName, "Error parsing data from page: "+err.Error())
		return
	}

	userId, userName, userInSystem := checkUserInSystem(deviceName, data.Data)
	var responseData ZapsiResponseData
	if !userInSystem {
		responseData.Data = "nok"
		responseData.UserId = userId
		responseData.UserName = userName
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(deviceName, "Check user finished, user not in system")
		return
	}
	responseData.Data = "ok"
	responseData.UserId = userId
	responseData.UserName = userName
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(deviceName, "Check user finished, everything ok")
}

func checkUserInSystem(deviceName string, user string) (int, string, bool) {
	logInfo(deviceName, "Checking user "+user)
	//TODO: check user in Zapsi
	if user == "12345" {
		return 23, "Brad Pitt", true
	}
	return 0, "", false
}

func endOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ipAddress := strings.Split(request.Host, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "End order in Zapsi called")
	var data OrderData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError(deviceName, "Error parsing data from page: "+err.Error())
		return
	}
	logInfo(deviceName, "Order: "+data.Order+"; userId:"+data.UserId+"; deviceId: "+data.DeviceId)
	//TODO: EndOrderInZapsi(data)
	var responseData ZapsiResponseData
	responseData.Data = "ok"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(deviceName, "End order in Zapsi finished, everything ok")
}

func getIdles(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ipAddress := strings.Split(request.Host, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Get idles from Zapsi called")
	var idles []Idle
	for i := 1; i <= 30; i++ {
		var idle Idle
		idle.Id = i
		idle.Name = "Prostoj c. " + strconv.Itoa(i)
		idles = append(idles, idle)
	}
	//TODO: get idles from database
	var responseData Idles
	responseData.Data = idles
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(deviceName, "Get idles from Zapsi finished, everything ok")
}

func endIdle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ipAddress := strings.Split(request.Host, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "End idle in Zapsi called")
	var data IdleData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError(deviceName, "Error parsing data from page: "+err.Error())
		return
	}
	logInfo(deviceName, "Order: "+data.Order+"; idleId: "+data.IdleId+"; userId: "+data.UserId+"; deviceId: "+data.DeviceId)
	//TODO: ENdIdleInZapsi(data)
	var responseData ZapsiResponseData
	responseData.Data = "ok"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(deviceName, "End idle in Zapsi finished, everything ok")
}

func startIdle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ipAddress := strings.Split(request.Host, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Start idle in Zapsi called")
	var data IdleData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError(deviceName, "Error parsing data from page: "+err.Error())
		return
	}
	logInfo(deviceName, "Order: "+data.Order+"; idleId: "+data.IdleId+"; userId: "+data.UserId+"; deviceId: "+data.DeviceId)
	//TODO: StartIdleInZapsi(data)
	var responseData ZapsiResponseData
	responseData.Data = "ok"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(deviceName, "Start idle in Zapsi finished, everything ok")
}
