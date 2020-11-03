package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"github.com/kardianos/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

const version = "2020.4.2.3"
const serviceName = "SK Label Cutting Webservice"
const serviceDescription = "Web Service for terminals for cutting workplaces"
const zapsiDatabaseConnection = "zapsi_uzivatel:zapsi@tcp(zapsidatabase:3306)/zapsi2?charset=utf8mb4&parseTime=True&loc=Local"
const skLabelDatabaseConnection = "sqlserver://zapsi:DSgEEmPNxCwgTJjsd2uR@10.3.1.3:1433?database=K2_SKLABEL"

var devicesMap map[string]string

type program struct{}

func (p *program) Start(service.Service) error {
	logInfo("MAIN", serviceName+" ["+version+"] started")
	go p.run()
	return nil
}

func (p *program) Stop(service.Service) error {
	logInfo("MAIN", serviceName+" ["+version+"] stopped")
	return nil
}

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

func (p *program) run() {
	router := httprouter.New()
	timer := sse.New()

	router.Handler("GET", "/time", timer)

	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/html/*filepath", http.Dir("html"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))

	router.GET("/", origin)
	router.GET("/cutting_end", cuttingEnd)
	router.GET("/home", home)
	router.GET("/idle_running", idleRunning)
	router.GET("/idle_select", idleSelect)
	router.GET("/user_change", userChange)
	router.GET("/user_break", userBreak)

	router.POST("/check_order", checkOrderInK2)
	router.POST("/check_user", checkUser)
	router.POST("/get_idles", getIdles)
	router.POST("/save_code", saveDataToK2)
	router.POST("/get_k2Pcs", getPcsFromK2)
	router.POST("/create_order", createOrder)
	router.POST("/create_idle", createIdle)
	router.POST("/end_idle", endIdle)
	router.POST("/end_order", endOrder)
	go streamTime(timer)
	go updateDeviceMap()
	err := http.ListenAndServe(":80", router)
	if err != nil {
		logError("MAIN", "Problem starting service: "+err.Error())
		os.Exit(-1)
	}
	logInfo("MAIN", serviceName+" ["+version+"] running")
}

func updateDeviceMap() {
	devicesMap = make(map[string]string)
	for {
		db, err := gorm.Open(mysql.Open(zapsiDatabaseConnection), &gorm.Config{})
		sqlDB, _ := db.DB()
		if err != nil {
			logError("MAIN", "Problem opening database: "+err.Error())
		}

		var devices []Device
		db.Where("DeviceType = 100").Find(&devices)
		for _, device := range devices {
			devicesMap[device.IPAddress] = device.Name
		}
		sqlDB.Close()
		time.Sleep(60 * time.Second)
	}
}

func streamTime(streamer *sse.Streamer) {
	logInfo("MAIN", "Streaming time process started")
	for {
		streamer.SendString("", "time", time.Now().Format("15:04"))
		time.Sleep(1 * time.Second)
	}
}
