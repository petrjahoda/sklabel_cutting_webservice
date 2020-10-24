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
	Result   string
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
	OrderBarcode string
	DeviceId     string
	UserId       string
	Pcs          string
}

type IdleData struct {
	OrderBarcode string
	IdleId       string
	DeviceId     string
	UserId       string
}

func checkIfUserIsLoggedForTerminalId(deviceName string, terminalId int) (int, string, bool) {
	logInfo(deviceName, "Checking if user is logged for terminal id: "+strconv.Itoa(terminalId))
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(deviceName, "Problem opening database: "+err.Error())
		return 0, "", false
	}

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
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(deviceName, "Problem opening database: "+err.Error())
		return 0, "", false
	}

	var device Device
	db.Where("IpAddress = ?", ipAddress).Where("DeviceType = 100").Find(&device)
	if len(device.Name) > 0 {
		logInfo(deviceName, "This ip address has assigned terminal: "+device.Name)
		var workplace Workplace
		db.Where("DeviceID = ?", device.OID).Find(&workplace)
		return device.OID, workplace.Code, true
	} else {
		logInfo(deviceName, "This ip address has not assigned any terminal")
		return 0, "", false
	}
}

func createOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ipAddress := strings.Split(request.RemoteAddr, ":")
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
	logInfo(deviceName, "OrderBarcode: "+data.OrderBarcode+"; userId:"+data.UserId+"; deviceId: "+data.DeviceId)
	actualWorkshiftId := GetActualWorkshiftId(deviceName, data.DeviceId)
	if actualWorkshiftId == 0 {
		logError(deviceName, "Problem getting workshift id")
		var responseData ZapsiResponseData
		responseData.Data = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		return
	}
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(deviceName, "Problem opening database: "+err.Error())
		var responseData ZapsiResponseData
		responseData.Data = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		return
	}
	var order Order
	db.Where("Barcode = ?", data.OrderBarcode).Find(&order)
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
	var workplace Workplace
	db.Where("DeviceId = ?", deviceIdInt).Find(&workplace)
	var workplaceMode WorkplaceMode
	db.Where("WorkplaceModeTypeID = 5").Where("WorkplaceID = ?", workplace.OID).Find(&workplaceMode)

	var terminalInputOrder TerminalInputOrder
	terminalInputOrder.DTS = time.Now()
	terminalInputOrder.OrderID = order.OID
	terminalInputOrder.DeviceID = deviceIdInt
	terminalInputOrder.Interval = 0
	terminalInputOrder.Count = 0
	terminalInputOrder.Fail = 0
	terminalInputOrder.AverageCycle = 0
	terminalInputOrder.WorkerCount = 0
	terminalInputOrder.WorkplaceModeID = workplaceMode.OID
	terminalInputOrder.WorkshiftID = actualWorkshiftId
	if userIdInt != 0 {
		terminalInputOrder.UserID = sql.NullInt32{
			Int32: int32(userIdInt),
			Valid: true,
		}
	}
	db.Save(&terminalInputOrder)
	var terminalInputLogin TerminalInputLogin
	db.Where("DeviceID = ?", data.DeviceId).Where("UserID = ?", data.UserId).Where("DTE is NULL").Find(&terminalInputLogin)
	if terminalInputLogin.OID > 0 {
		logInfo(deviceName, "User already logged in terminal_input_login")
	} else if userIdInt != 0 {
		logInfo(deviceName, "Logging user to terminal_input_login")
		var newTerminalInputLogin TerminalInputLogin
		newTerminalInputLogin.DTS = time.Now()
		newTerminalInputLogin.DeviceID = deviceIdInt
		newTerminalInputLogin.UserID = userIdInt
		db.Save(&newTerminalInputLogin)
	}
	var responseData ZapsiResponseData
	responseData.Data = "ok"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(deviceName, "Create order in Zapsi finished, everything ok")
}

func GetActualWorkshiftId(deviceName string, deviceID string) int {
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(deviceName, "Problem opening database: "+err.Error())
		return 0
	}
	var workplace Workplace
	db.Where("DeviceID = ?", deviceID).Find(&workplace)
	var workplaceDivision WorkplaceDivision
	db.Where("OID = ?", workplace.WorkplaceDivisionID).Find(&workplaceDivision)
	var workShifts []Workshift
	db.Where("Active = 1").Find(&workShifts)
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
	ipAddress := strings.Split(request.RemoteAddr, ":")
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
	logInfo(deviceName, "Data before parsing: "+data.Data)
	updatedData := strings.ReplaceAll(data.Data, "SHIFT", "")
	updatedData = strings.ReplaceAll(updatedData, "ENTER", "")
	updatedData = strings.ReplaceAll(updatedData, "/R", "")
	logInfo(deviceName, "Data after parsing: "+updatedData)
	userId, userName, userInSystem := checkUserInZapsi(deviceName, updatedData)
	if !userInSystem {
		var responseData ZapsiResponseData
		responseData.Data = "nok"
		responseData.UserId = userId
		responseData.UserName = userName
		responseData.Result = updatedData
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo(deviceName, "Check user finished, user not in system")
		return
	}
	var responseData ZapsiResponseData
	responseData.Data = "ok"
	responseData.UserId = userId
	responseData.UserName = userName
	responseData.Result = updatedData
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(deviceName, "Check user finished, everything ok")
}

func checkUserInZapsi(deviceName string, userRfid string) (int, string, bool) {
	logInfo(deviceName, "Checking user "+userRfid)
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
	}
	var user User
	db.Where("Rfid = ?", userRfid).Find(&user)
	if user.OID > 0 {
		return user.OID, user.FirstName + " " + user.Name, true
	}
	return 0, "", false
}

func endOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ipAddress := strings.Split(request.RemoteAddr, ":")
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
	logInfo(deviceName, "OrderBarcode: "+data.OrderBarcode+"; userId:"+data.UserId+"; deviceId: "+data.DeviceId)
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(deviceName, "Problem opening database: "+err.Error())
		responseData.Data = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	}
	pcsToInsert, err := strconv.Atoi(data.Pcs)
	if err != nil {
		logError(deviceName, "Problem parsing count: "+err.Error())
		pcsToInsert = 0
	}

	var runningOrder TerminalInputOrder
	db.Where("DeviceID = ?", data.DeviceId).Where("DTE is NULL").Find(&runningOrder)
	runningOrder.DTE = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	runningOrder.Interval = float32(time.Since(runningOrder.DTS).Minutes())
	runningOrder.Count = pcsToInsert
	runningOrder.AverageCycle = float32(time.Since(runningOrder.DTS).Minutes()) / float32(pcsToInsert)
	db.Save(&runningOrder)
	var terminalInputLogin TerminalInputLogin
	db.Where("DeviceID = ? ", data.DeviceId).Where("DTE is NULL").Find(&terminalInputLogin)
	if terminalInputLogin.OID > 0 {
		logInfo(deviceName, "Closing terminal_input_login")
		terminalInputLogin.DTE = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		terminalInputLogin.Interval = float32(time.Since(terminalInputLogin.DTS).Minutes())
		db.Save(&terminalInputLogin)
	} else {
		logError(deviceName, "No terminal_input_login found when closing order")
	}
	responseData.Data = "ok"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo(deviceName, "End order in Zapsi finished, everything ok")
}

func getIdles(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ipAddress := strings.Split(request.RemoteAddr, ":")
	deviceName := devicesMap[ipAddress[0]]
	if len(deviceName) == 0 {
		deviceName = ipAddress[0]
	}
	logInfo(deviceName, "Get idles from Zapsi called")
	var dataIdles []DataIdle
	var responseData DataIdles
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(deviceName, "Problem opening database: "+err.Error())
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	}

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
	ipAddress := strings.Split(request.RemoteAddr, ":")
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
	logInfo(deviceName, "OrderBarcode: "+data.OrderBarcode+"; idleId: "+data.IdleId+"; userId: "+data.UserId+"; deviceId: "+data.DeviceId)
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
	}

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
	ipAddress := strings.Split(request.RemoteAddr, ":")
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
	logInfo(deviceName, "OrderBarcode: "+data.OrderBarcode+"; idleId: "+data.IdleId+"; userId: "+data.UserId+"; deviceId: "+data.DeviceId)
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
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
	}

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

func checkOrderInZapsi(deviceName string, skZapsiVp SkZapsiVp) {
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(deviceName, "Problem opening database: "+err.Error())
		return
	}

	var order Order
	db.Where("Barcode = ?", skZapsiVp.VPexp).Find(&order)
	if order.OID > 0 {
		logInfo(deviceName, "OrderBarcode with barcode "+skZapsiVp.VPexp+" already exists")
	} else {
		logInfo(deviceName, "OrderBarcode with barcode "+skZapsiVp.VPexp+" does not exists, creating")
		var newOrder Order
		newOrder.Name = skZapsiVp.VP
		newOrder.Barcode = skZapsiVp.VPexp
		newOrder.ProductID = 100
		newOrder.OrderStatusID = 1
		newOrder.CountRequested = 0
		db.Save(&newOrder)
	}
}

func GetUserLoginFor(deviceName string, userId string) string {
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(deviceName, "Problem opening database: "+err.Error())
		return ""
	}

	var user User
	db.Where("OID = ?", userId).Find(&user)
	return user.Login
}

func GetIdleBarcodeFor(deviceName string, idleId string) string {
	db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(deviceName, "Problem opening database: "+err.Error())
		return ""
	}

	var idle Idle
	db.Where("OID = ?", idleId).Find(&idle)
	return idle.Barcode
}
