package main

import (
	"context"
	"flag"
	"github.com/praveenprem/testbed-slack-bot/config"
	"github.com/praveenprem/testbed-slack-bot/slack"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	VERSION = "0.0.1-SNAPSHOT"
)

func main() {
	cfgPath := flag.String("config-path", "", "Configuration file path")
	flag.Parse()

	if *cfgPath == "" {
		log.Fatalln("missing arg --config-path")
	}

	log.Println("Application startup")
	log.Printf("Version: %s", VERSION)

	var cfg config.Config
	if err := cfg.Load(*cfgPath); err != nil {
		log.Fatalf("%#v", err.Error())
	}

	slack.CFG = cfg.Slack
	serverConfig := cfg
	server := serverConfig.Start()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("Server up.....")
	<-shutdown

	log.Println("Starting cleanup")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Server shutdown completed")

}
