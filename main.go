package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"github.com/kardianos/service"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"
)

type SaveToK2Data struct {
	Data string
}

type ResponseData struct {
	Data     string
	UserId   int
	UserName string
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

type Idle struct {
	Id   int
	Name string
}

type Idles struct {
	Data []Idle
}

const version = "2020.4.1.2"
const serviceName = "SK Label Cutting Webservice"
const serviceDescription = "Web Service for terminals for cutting workplaces"
const zapsiDatabaseConnection = "user=postgres password=Zps05..... dbname=version3 host=database port=5432 sslmode=disable"
const skLabelDatabaseConnection = "user=postgres password=Zps05..... dbname=version3 host=database port=5432 sslmode=disable"

type program struct{}

func main() {
	logInfo("MAIN", serviceName+" ["+version+"] starting...")
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
	err = s.Run()
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
}

func (p *program) Start(service.Service) error {
	logInfo("MAIN", serviceName+" ["+version+"] started")
	go p.run()
	return nil
}

func (p *program) Stop(service.Service) error {
	logInfo("MAIN", serviceName+" ["+version+"] stopped")
	return nil
}

func (p *program) run() {
	router := httprouter.New()
	timer := sse.New()

	router.Handler("GET", "/time", timer)

	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/html/*filepath", http.Dir("html"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))

	router.GET("/", origin)
	router.GET("/user_error", userError)
	router.GET("/cutting_end", cuttingEnd)
	router.GET("/home", home)
	router.GET("/idle_running", idleRunning)
	router.GET("/idle_select", idleSelect)
	router.GET("/user_change", userChange)
	router.GET("/user_break", userBreak)

	router.POST("/check_order", checkOrder)
	router.POST("/check_user", checkUser)
	router.POST("/get_idles", getIdles)
	router.POST("/save_code", saveCode)
	router.POST("/get_k2Pcs", getK2Pcs)
	router.POST("/start_order", startOrder)
	router.POST("/start_idle", startIdle)
	router.POST("/end_idle", endIdle)
	router.POST("/end_order", endOrder)

	go streamTime(timer)
	err := http.ListenAndServe(":80", router)
	if err != nil {
		logError("MAIN", "Problem starting service: "+err.Error())
		os.Exit(-1)
	}
	logInfo("MAIN", serviceName+" ["+version+"] running")
}

func getK2Pcs(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Get K2 Pcs", "Getting K2 pcs called")
	var data SaveToK2Data
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("Get K2 Pcs", err.Error())
		return
	}
	logInfo("Get K2 Pcs", "Getting K2 pcs for order "+data.Data)
	//TODO: Get Pcs from K2
	var responseData ResponseData
	responseData.Data = "23"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("Get K2 Pcs", "Getting K2 pcs finished")
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

func checkUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Check User", "Checking user called")
	var data SaveToK2Data
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("Check User", err.Error())
		return
	}

	userId, userName, userInSystem := checkUserInSystem(data.Data)
	var responseData ResponseData
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

func userBreak(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	logInfo("User Break", "Page loading...")
	var data HomePageData
	data.Time = time.Now().Format("15:04:05")
	tmpl := template.Must(template.ParseFiles("./html/user_break.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("User Break", "Page loaded")
}

func userChange(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	logInfo("Login", "Page loading...")
	var data HomePageData
	data.Time = time.Now().Format("15:04:05")
	tmpl := template.Must(template.ParseFiles("./html/user_change.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("Login", "Page loaded")
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

func saveCode(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Save Code", "Saving code called")
	var data SaveToK2Data
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", err.Error())
		return
	}
	logInfo("Save Code", "Saving code "+data.Data)
	//TODO: SaveCodeToK2
	logInfo("Save Code", "Saving code finished")

}

func checkOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("Check Order", "Checking order called")
	var data SaveToK2Data
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", err.Error())
		return
	}

	orderIsInSystem := checkOrderInSystem(data.Data)
	var responseData ResponseData
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

func streamTime(streamer *sse.Streamer) {
	logInfo("SSE", "Streaming time process started")
	for {
		streamer.SendString("", "time", time.Now().Format("15:04:05"))
		time.Sleep(1 * time.Second)
	}
}
