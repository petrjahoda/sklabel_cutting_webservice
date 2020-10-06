package main

import (
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ZapsiResponseData struct {
	Data     string
	UserId   int
	UserName string
}

type ZapsiData struct {
	Data string
}

type DataIdle struct {
	Id   int
	Name string
}

type DataIdles struct {
	Data []DataIdle
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
	logInfo(deviceName, "Checking if user is logged for terminal id: "+strconv.Itoa(terminalId))
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(deviceName, "Problem opening database: "+err.Error())
		return 0, "", false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var terminalInputLogin TerminalInputLogin
	db.Where("DeviceID = ?", terminalId).Where("DTE is NULL").Find(&terminalInputLogin)
	var user User
	db.Where("Oid = ?", terminalInputLogin.UserID).Find(&user)
	if user.OID > 0 {
		logInfo(deviceName, "Logged user: "+user.FirstName+" "+user.Name)
		return user.OID, user.FirstName + " " + user.Name, true
	} else {
		logInfo(deviceName, "This device has no logged user")
		return 0, "", false
	}
}

func checkIpAddress(deviceName string, ipAddress string) (int, string, bool) {
	logInfo(deviceName, "Checking ip address: "+ipAddress)
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(deviceName, "Problem opening database: "+err.Error())
		return 0, "", false
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var device Device
	db.Where("IpAddress = ?", ipAddress).Where("DeviceType = 100").Find(&device)
	if len(device.Name) > 0 {
		logInfo(deviceName, "This ip address has assigned terminal: "+device.Name)
		return device.OID, device.Name, true
	} else {
		logInfo(deviceName, "This ip address has not assigned any terminal")
		return 0, "", false
	}
}

func createOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ipAddress := strings.Split(request.Host, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Create order in Zapsi called")
	var data OrderData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError(deviceName, "Error parsing data from page: "+err.Error())
		var responseData ZapsiResponseData
		responseData.Data = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		return
	}
	logInfo(deviceName, "Order: "+data.Order+"; userId:"+data.UserId+"; deviceId: "+data.DeviceId)
	actualWorkshiftId := GetActualWorkshiftId(deviceName, data.DeviceId)
	if actualWorkshiftId == 0 {
		logError(deviceName, "Problem getting workshift id")
		var responseData ZapsiResponseData
		responseData.Data = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		return
	}
	createOrderIfNotPresent(deviceName, data.Order)
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(deviceName, "Problem opening database: "+err.Error())
		var responseData ZapsiResponseData
		responseData.Data = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		return
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var order Order
	db.Where("Barcode = ?", data.Order).Find(&order)
	userIdInt, err := strconv.Atoi(data.UserId)
	if err != nil {
		logError(deviceName, "Problem parsing userid "+data.UserId+": "+err.Error())

	}
	deviceIdInt, err := strconv.Atoi(data.DeviceId)
	if err != nil {
		logError(deviceName, "Problem parsing deviceid "+data.DeviceId+": "+err.Error())
		var responseData ZapsiResponseData
		responseData.Data = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		return
	}
	var terminalInputOrder TerminalInputOrder
	terminalInputOrder.DTS = time.Now()
	terminalInputOrder.OrderID = order.OID
	terminalInputOrder.DeviceID = deviceIdInt
	terminalInputOrder.Interval = 0
	terminalInputOrder.Count = 0
	terminalInputOrder.Fail = 0
	terminalInputOrder.AverageCycle = 0
	terminalInputOrder.WorkerCount = 0
	terminalInputOrder.WorkplaceModeID = 1 //TODO: upravit tady spravne
	terminalInputOrder.WorkshiftID = actualWorkshiftId
	if userIdInt != 0 {
		terminalInputOrder.UserID = sql.NullInt32{
			Int32: int32(userIdInt),
			Valid: true,
		}
	}
	db.Save(&terminalInputOrder)
	var responseData ZapsiResponseData
	responseData.Data = "ok"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(deviceName, "Create order in Zapsi finished, everything ok")
}

func createOrderIfNotPresent(deviceName string, orderBarcode string) {
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(deviceName, "Problem opening database: "+err.Error())
		return
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var order Order
	db.Where("Barcode = ?", orderBarcode).Find(&order)
	if order.OID > 0 {
		logInfo(deviceName, "Order already exists")
		return
	} else {
		logInfo(deviceName, "Order does not exist, creating")
		var newOrder Order
		newOrder.Name = orderBarcode //TODO: upravit tady spravne
		newOrder.Barcode = orderBarcode
		newOrder.ProductID = 101   //TODO: upravit tady spravne
		newOrder.OrderStatusID = 1 //TODO: upravit tady spravne
		newOrder.WorkplaceID = 1   //TODO: upravit tady spravne
		db.Save(&newOrder)
	}
}

func GetActualWorkshiftId(deviceName string, deviceID string) int {
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(deviceName, "Problem opening database: "+err.Error())
		return 0
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var workplace Workplace
	db.Where("DeviceId = ?", deviceID).Find(&workplace)
	var workplaceDivision WorkplaceDivision
	db.Where("OID = ?", workplace.WorkplaceDivisionID).Find(&workplaceDivision)
	var workShifts []Workshift
	db.Where("Active = 1").Where("WorkplaceDivisionID = ?", workplaceDivision.OID).Find(&workShifts)
	logInfo(deviceName, "Found "+strconv.Itoa(len(workShifts))+" active workshifts")
	if len(workShifts) > 0 {
		logInfo(deviceName, "Finding proper workshift")
		for _, workshift := range workShifts {
			startDate := strings.Split(workshift.WorkshiftStart, ":")
			if len(startDate) > 2 {
				startHour, err := strconv.Atoi(startDate[0])
				startMinute, err := strconv.Atoi(startDate[1])
				if err != nil {
					logError(deviceName, "WorkshiftStart does not have proper format HH:MM:SS "+workshift.WorkshiftStart)
					continue
				}
				now := time.Now()
				start := time.Time{}
				if startHour > now.Hour() {
					start = time.Date(now.Year(), now.Month(), now.Day()-1, startHour, startMinute, 00, 0, time.Local)
				} else {
					start = time.Date(now.Year(), now.Month(), now.Day(), startHour, startMinute, 00, 0, time.Local)
				}
				end := start.Add(time.Duration(workshift.WorkshiftLenght) * time.Minute)
				if start.Before(now) && now.Before(end) {
					logInfo(deviceName, "Actual workshiftid: "+strconv.Itoa(workshift.OID))
					return workshift.OID
				}
			} else {
				logError(deviceName, "WorkshiftStart does not have proper format HH:MM:SS "+workshift.WorkshiftStart)
			}
		}
		return 0
	}
	return 0
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
		var responseData ZapsiResponseData
		responseData.Data = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		return
	}

	userId, userName, userInSystem := checkUserInZapsi(deviceName, data.Data)
	if !userInSystem {
		var responseData ZapsiResponseData
		responseData.Data = "nok"
		responseData.UserId = userId
		responseData.UserName = userName
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(deviceName, "Check user finished, user not in system")
		return
	}
	var responseData ZapsiResponseData
	responseData.Data = "ok"
	responseData.UserId = userId
	responseData.UserName = userName
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(deviceName, "Check user finished, everything ok")
}

func checkUserInZapsi(deviceName string, userRfid string) (int, string, bool) {
	logInfo(deviceName, "Checking user "+userRfid)
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var user User
	db.Where("Rfid = ?", userRfid).Find(&user)
	if user.OID > 0 {
		return user.OID, user.FirstName + " " + user.Name, true
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
	var responseData ZapsiResponseData
	if err != nil {
		logError(deviceName, "Error parsing data from page: "+err.Error())
		responseData.Data = "nok"
		writer.Header().Set("Content-Type", "application/json")
		return
	}
	logInfo(deviceName, "Order: "+data.Order+"; userId:"+data.UserId+"; deviceId: "+data.DeviceId)
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(deviceName, "Problem opening database: "+err.Error())
		responseData.Data = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var runningOrder TerminalInputOrder
	db.Where("DeviceID = ?", data.DeviceId).Where("DTE is NULL").Find(&runningOrder)
	runningOrder.DTE = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	runningOrder.Interval = float32(time.Since(runningOrder.DTS).Minutes())
	db.Save(&runningOrder)
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
	var dataIdles []DataIdle
	var responseData DataIdles
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError(deviceName, "Problem opening database: "+err.Error())
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var idles []Idle
	db.Find(&idles)
	for _, idle := range idles {
		var dataIdle DataIdle
		dataIdle.Id = idle.OID
		dataIdle.Name = idle.Name
		dataIdles = append(dataIdles, dataIdle)
	}
	responseData.Data = dataIdles
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
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var runningIdle TerminalInputIdle
	db.Where("IdleID = ?", data.IdleId).Where("UserID = ? ", data.UserId).Where("DeviceID = ?", data.DeviceId).Where("DTE is NULL").Find(&runningIdle)
	runningIdle.DTE = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	runningIdle.Interval = float32(time.Since(runningIdle.DTS).Minutes())
	db.Save(&runningIdle)
	var responseData ZapsiResponseData
	responseData.Data = "ok"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(deviceName, "End idle in Zapsi finished, everything ok")
}

func createIdle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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
	deviceIdInt, err := strconv.Atoi(data.DeviceId)
	userIdInt, err := strconv.Atoi(data.UserId)
	idleIdInt, err := strconv.Atoi(data.IdleId)
	if err != nil {
		logError(deviceName, "Problem parsing data from user: "+err.Error())
		var responseData ZapsiResponseData
		responseData.Data = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		return
	}
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
	}
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	var terminalInputIdle TerminalInputIdle
	terminalInputIdle.DTS = time.Now()
	terminalInputIdle.IdleID = idleIdInt
	terminalInputIdle.UserID = userIdInt
	terminalInputIdle.Interval = 0
	terminalInputIdle.DeviceID = deviceIdInt
	db.Save(&terminalInputIdle)
	var responseData ZapsiResponseData
	responseData.Data = "ok"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(deviceName, "Start idle in Zapsi finished, everything ok")
}
