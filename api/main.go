package main

import (
	"greeny/routers"
	"greeny/utils"
	"log"
	"net/http"
	"time"
)

var configPath = "conf/conf.yml"

func main() {
	err := utils.SetConfigPath(configPath)
	if err != nil {
		log.Println(err)
		return
	}
	conf, err := utils.GetConfig()
	if err != nil {
		return
	}

	srv := &http.Server{
		Handler:      routers.Router(),
		Addr:         conf.Server.Host + ":" + conf.Server.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
