package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"github.com/kardianos/service"
	"net/http"
	"os"
	"time"
)

const version = "2020.4.1.2"
const serviceName = "SK Label Cutting Webservice"
const serviceDescription = "Web Service for terminals for cutting workplaces"
const zapsiDatabaseConnection = "user=postgres password=Zps05..... dbname=version3 host=database port=5432 sslmode=disable"
const skLabelDatabaseConnection = "user=postgres password=Zps05..... dbname=version3 host=database port=5432 sslmode=disable"

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

func streamTime(streamer *sse.Streamer) {
	logInfo("SSE", "Streaming time process started")
	for {
		streamer.SendString("", "time", time.Now().Format("15:04:05"))
		time.Sleep(1 * time.Second)
	}
}
