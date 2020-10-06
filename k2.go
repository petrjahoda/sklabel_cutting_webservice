package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type K2ResponseData struct {
	Data     string
	UserId   int
	UserName string
}

type SaveToK2Data struct {
	Data string
}

func saveCode(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("Save Code", "Saving code called")
	var data SaveToK2Data
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", err.Error())
		return
	}
	logInfo("Save Code", "Saving code "+data.Data)
	var responseData K2ResponseData
	responseData.Data = "ok"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	//TODO: SaveCodeToK2
	logInfo("Save Code", "Saving code finished")

}

func checkOrder(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("Check Order", "Checking order called")
	var data SaveToK2Data
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", err.Error())
		return
	}
	orderIsInSystem := checkOrderInSystem(data.Data)
	var responseData K2ResponseData
	if !orderIsInSystem {
		responseData.Data = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("Check Order", "Checking order finished")
		return
	}
	responseData.Data = "ok"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("Check Order", "Checking order finished")
}

func checkOrderInSystem(order string) bool {
	logInfo("Check Order In System", "Checking order "+order)
	logInfo("Check Order In System", "Order checked ")
	//TODO: check order in K2
	if order == "12345" {
		return true
	}
	return false
}

func getK2Pcs(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("Get K2 Pcs", "Getting K2 pcs called")
	var data SaveToK2Data
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("Get K2 Pcs", err.Error())
		return
	}
	logInfo("Get K2 Pcs", "Getting K2 pcs for order "+data.Data)
	//TODO: Get Pcs from K2
	var responseData K2ResponseData
	responseData.Data = "23"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("Get K2 Pcs", "Getting K2 pcs finished")
}
