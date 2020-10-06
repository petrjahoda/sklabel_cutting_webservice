package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

type K2ResponseData struct {
	Data     string
	UserId   int
	UserName string
}

type SaveToK2Data struct {
	Data string
}

func checkUserInK2(deviceName string, user string) (int, string, bool) {
	logInfo(deviceName, "Checking user "+user)
	//TODO: check user in Zapsi
	if user == "12345" {
		return 23, "Brad Pitt", true
	}
	return 0, "", false
}

func saveDataToK2(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ipAddress := strings.Split(request.Host, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Saving data to K2 called")
	var data SaveToK2Data
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError(deviceName, "Error parsing data from page: "+err.Error())
		return
	}
	logInfo(deviceName, "Data: "+data.Data)
	var responseData K2ResponseData
	responseData.Data = "ok"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	//TODO: SaveCodeToK2
	logInfo(deviceName, "Saving data to K2 finished, everything ok")

}

func checkOrderInK2(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ipAddress := strings.Split(request.Host, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Check order in K2 called")
	var data SaveToK2Data
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError(deviceName, "Error parsing data from page: "+err.Error())
		return
	}
	logInfo(deviceName, "Data: "+data.Data)
	orderIsInSystem := checkOrderInSystem(deviceName, data.Data)
	var responseData K2ResponseData
	if !orderIsInSystem {
		responseData.Data = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(deviceName, "Checking order in K2 finished, order not in system ")
		return
	}
	responseData.Data = "ok"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(deviceName, "Check order in K2 finished, everything ok")
}

func checkOrderInSystem(deviceName string, order string) bool {
	logInfo(deviceName, "Checking order "+order)
	//TODO: check order in K2
	if order == "12345" {
		return true
	}
	return false
}

func getPcsFromK2(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ipAddress := strings.Split(request.Host, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Get pcs from K2 called")
	var data SaveToK2Data
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError(deviceName, "Error parsing data from page: "+err.Error())
		return
	}
	logInfo(deviceName, "Data: "+data.Data)
	//TODO: Get Pcs from K2
	var responseData K2ResponseData
	responseData.Data = "23"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(deviceName, "Get pcs from K2 finished, everything ok")
}
