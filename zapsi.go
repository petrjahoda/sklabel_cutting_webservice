package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
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

type StartOrderData struct {
	Order    string
	DeviceId string
	UserId   string
}

type StartIdleData struct {
	Order    string
	IdleId   string
	DeviceId string
	UserId   string
}

func checkIfUserIsLoggedForTerminalId(terminalId int) (int, string, bool) {
	//TODO: check if user is logged in terminal input user (nebo jak)
	logInfo("Origin", "Checking if user is logged for terminal id: "+strconv.Itoa(terminalId))
	logInfo("Origin", "User is logged")
	return 23, "Petr Jahoda", true
}

func checkIpAddress(ipAddress string) (int, string, bool) {
	//TODO: check assigned terminal to this ip
	logInfo("Origin", "Checking ip address: "+ipAddress)
	logInfo("Origin", "Ip address assigned to terminal")
	return 1, "E235", true
}

func startOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Start Order", "Saving order in Zapsi called")
	var data StartOrderData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", err.Error())
		return
	}
	logInfo("Start Order", "Order: "+data.Order+"; userId:"+data.UserId+"; deviceId: "+data.DeviceId)
	//TODO: StartOrderInZapsi(data)
	logInfo("Start Order", "Saving order in Zapsi finished")
}

func checkUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Check User", "Checking user called")
	var data ZapsiData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("Check User", err.Error())
		return
	}

	userId, userName, userInSystem := checkUserInSystem(data.Data)
	var responseData ZapsiResponseData
	if !userInSystem {
		responseData.Data = "nok"
		responseData.UserId = userId
		responseData.UserName = userName
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check User", "Checking user finished")
		return
	}
	responseData.Data = "ok"
	responseData.UserId = userId
	responseData.UserName = userName
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("Check User", "Checking user finished")
}

func checkUserInSystem(user string) (int, string, bool) {
	logInfo("Check User In System", "Checking user "+user)
	logInfo("Check User In System", "Order user ")
	//TODO: check user in Zapsi
	if user == "12345" {
		return 23, "Brad Pitt", true
	}
	return 0, "", false
}

func endOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("End Order", "Ending order in Zapsi called")
	var data StartOrderData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("End Order", err.Error())
		return
	}
	logInfo("End Order", "Order: "+data.Order+"; userId:"+data.UserId+"; deviceId: "+data.DeviceId)
	//TODO: EndOrderInZapsi(data)
	logInfo("End Order", "Ending order in Zapsi finished")
}

func getIdles(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Get Idles", "Getting idles called")
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
	logInfo("Get Idles User", "Getting idles finished")
}

func endIdle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("End Idle", "Ending idle in Zapsi called")
	var data StartIdleData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("End Idle", err.Error())
		return
	}
	logInfo("End Idle", "Order: "+data.Order+"; idleId: "+data.IdleId+"; userId: "+data.UserId+"; deviceId: "+data.DeviceId)
	//TODO: ENdIdleInZapsi(data)
	logInfo("End Idle", "Ending idle in Zapsi finished")
}

func startIdle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Start Idle", "Saving idle in Zapsi called")
	var data StartIdleData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("Start Idle", err.Error())
		return
	}
	logInfo("Start Idle", "Order: "+data.Order+"; idleId: "+data.IdleId+"; userId: "+data.UserId+"; deviceId: "+data.DeviceId)
	//TODO: StartIdleInZapsi(data)
	logInfo("Start Idle", "Saving idle in Zapsi finished")
}
