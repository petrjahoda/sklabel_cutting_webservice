package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"github.com/kardianos/service"
	"net/http"
	"os"
	"time"
)

type SaveToK2Data struct {
	Data string
}

type ResponseData struct {
	Data string
}

type StartOrderData struct {
	Order    string
	DeviceId string
	UserId   string
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
	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/html/*filepath", http.Dir("html"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))

	router.Handler("GET", "/time", timer)

	router.GET("/", origin)
	router.GET("/user_error", userError)
	router.GET("/entry_pcs", entryPcs)
	router.GET("/home", home)
	router.GET("/idle_running", idleRunning)
	router.GET("/idle_select", idleSelect)
	router.GET("/login", login)

	router.POST("/check_order", checkOrder)
	router.POST("/save_code", saveCode)
	router.POST("/start_order", startOrder)

	go streamTime(timer)
	err := http.ListenAndServe(":80", router)
	if err != nil {
		logError("MAIN", "Problem starting service: "+err.Error())
		os.Exit(-1)
	}
	logInfo("MAIN", serviceName+" ["+version+"] running")
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
	//TODO: SaveToZapsi(data)
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
	//TODO: SaveToK2(data.Data)
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
	//TODO: check order
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
