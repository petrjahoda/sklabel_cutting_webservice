package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"github.com/kardianos/service"
	"net/http"
	"os"
	"time"
)

const version = "2020.3.3.28"
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

	router.GET("/", entry)
	router.GET("/order_scanning", orderScanning)
	router.GET("/order_scanning", orderError)
	router.GET("/order_scanning", userError)
	router.GET("/entry_pcs", entryPcs)
	router.GET("/home", home)
	router.GET("/idle_running", idleRunning)
	router.GET("/idle_select", idleSelect)
	router.GET("/login", login)
	
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
		location, err := time.LoadLocation("Europe/Prague")
		if err != nil {
			logError("MAIN", "Problem loading timezone for Europe/Prague")
		} else {
			streamer.SendString("", "time", time.Now().In(location).Format("15:04:05"))
			time.Sleep(1 * time.Second)
		}
	}
}
