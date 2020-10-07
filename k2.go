package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type K2ResponseData struct {
	Data     string
	UserId   int
	UserName string
}

type SaveToK2Data struct {
	Type          string
	Code          string
	WorkplaceCode string
	UserId        string
	OrderBarcode  string
	IdleId        string
	Pcs           string
}
type K2Data struct {
	Data string
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
		var responseData K2ResponseData
		responseData.Data = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		return
	}
	db, err := gorm.Open(sqlserver.Open(skLabelDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		var responseData K2ResponseData
		responseData.Data = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		return
	}
	userLogin := GetUserLoginFor(deviceName, data.UserId)
	idleBarcode := GetIdleBarcodeFor(deviceName, data.IdleId)
	dataToInsert := ""
	if data.Type == "order" {
		dataToInsert = "\\id_stroj{" + data.WorkplaceCode + "}\\id_osoby{" + userLogin + "}\\id_zakazky{" + data.OrderBarcode + "}\\id_krok{" + data.Code + "}\\id_operace{" + data.OrderBarcode + "}"
		if data.Code == "K302" {
			dataToInsert += "\\pocet_impulzu{" + data.Pcs + "}"
		}
	} else {
		dataToInsert = "\\id_stroj{" + data.WorkplaceCode + "}\\id_osoby{" + userLogin + "}\\id_zakazky{" + data.OrderBarcode + "}\\id_krok{" + data.Code + "}\\id_operace{" + data.OrderBarcode + "}\\duvod{" + idleBarcode + "}"
	}
	logInfo(deviceName, "K2 STRING: "+dataToInsert)
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var zapsiK2 ZapsiK2
	zapsiK2.Cas = time.Now()
	zapsiK2.Typ = 200
	zapsiK2.Data = dataToInsert
	zapsiK2.Zprac = 0
	//db.Save(&zapsiK2)
	var responseData K2ResponseData
	responseData.Data = "ok"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(deviceName, "Saving data to K2 finished, everything ok")

}

func checkOrderInK2(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ipAddress := strings.Split(request.Host, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Check order in K2 called")
	var data K2Data
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError(deviceName, "Error parsing data from page: "+err.Error())
		return
	}
	logInfo(deviceName, "Data: "+data.Data)
	skZapsiVP, orderIsInSystem := checkOrderInSystem(deviceName, data.Data)
	var responseData K2ResponseData
	if !orderIsInSystem {
		responseData.Data = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(deviceName, "Checking order in K2 finished, order not in K2 database ")
		return
	}
	checkOrderInZapsi(deviceName, skZapsiVP)
	responseData.Data = "ok"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(deviceName, "Check order in K2 finished, everything ok")
}

func checkOrderInSystem(deviceName string, order string) (SkZapsiVp, bool) {
	logInfo(deviceName, "Checking order in K2: "+order)
	db, err := gorm.Open(sqlserver.Open(skLabelDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return SkZapsiVp{}, false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var skZapsiVp SkZapsiVp
	db.Debug().Where("MaterialBM IS NOT NULL").Where("MaterialBL > 0").Where("VPexp = ?", order).Find(&skZapsiVp)
	if skZapsiVp.RID > 0 {
		return skZapsiVp, true
	} else {
		return skZapsiVp, false
	}
}

func getPcsFromK2(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ipAddress := strings.Split(request.Host, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Get pcs from K2 called")
	var data K2Data
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError(deviceName, "Error parsing data from page: "+err.Error())
		return
	}
	logInfo(deviceName, "Data: "+data.Data)

	db, err := gorm.Open(sqlserver.Open(skLabelDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		var responseData K2ResponseData
		responseData.Data = "0"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(deviceName, "Get pcs from K2 finished, cannot get data from K2")
		return
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var skZapsiVp SkZapsiVp
	db.Debug().Where("MaterialBM IS NOT NULL").Where("MaterialBL > 0").Where("VPexp = ?", data.Data).Find(&skZapsiVp)
	var responseData K2ResponseData
	if skZapsiVp.RID > 0 {
		data := strconv.FormatFloat(float64(skZapsiVp.MaterialBM), 'g', 1, 32)
		if err != nil {
			logError(deviceName, "Cannot parse MaterialBM: "+err.Error())
			responseData.Data = "0"
		} else {
			responseData.Data = data
		}
	} else {
		responseData.Data = "0"
	}

	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(deviceName, "Get pcs from K2 finished, everything ok")
}
