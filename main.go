package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/mcereal/botty/config"
	"github.com/mcereal/botty/cron"
	"github.com/mcereal/botty/router"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Load config and environment
	config.LoadConfiguration()
	error := godotenv.Load()
	if error != nil {
		log.Warn("Error loading .env file. Proceeding with default config. ")
	}

	// Start the Stale PR cron job
	cron.ScheduleCron()

	/*
	 Check for environment variables.
	 It they are not present port will be set to the config value from
	 the config.yml. if config.yml is not present then a default
	 value will be loaded from app_config
	*/

	// check for the "PORT" environment variable.
	port := os.Getenv("PORT")
	if port == "" {
		port = config.AppConfig.Server.Port
	}

	// check for the "ENVIRONMENT" environment variable.
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = config.AppConfig.Application.Environment
	}

	// Initialize gin and pass in the current environment
	r := router.InitializeRouter(env)

	// create http server
	srv := &http.Server{}

	/*
	 populate http server with environment defaults.
	 if the environment is development "localhost" is added
	 in the address. This is not neccesary, but prevents pop up
	 dialog every time you start the server.
	*/
	if env == "development" {
		log.Infof("Environment: %s", env)
		log.Infof("Listening on http://localhost:%s/ ", port)
		srv = &http.Server{
			Addr:    "localhost:" + port,
			Handler: r,
		}
	} else {
		address := fmt.Sprintf(":%s", port)
		log.Infof("Environment: %s", env)
		log.Infof("Listening on port: %s", port)
		srv = &http.Server{
			Addr:    address,
			Handler: r,
		}
	}

	/*
		Gracefully exit server when the operating system signal is interupted.
		A channel is created that takes in the the the os signal. Then a go routine
		which is a lightweight thread execution runs in the background. When it
		recieves the "quit" variable which is mapped too oc. Interrupt then it will
		close the server.
	*/

	// create a channel that takes in an interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// go routine to close server when channel recieves interupt
	go func() {
		<-quit
		log.Println("receive interrupt signal")
		if err := srv.Close(); err != nil {
			log.Fatal("Server Close:", err)
		}
	}()

	// start the server
	if err := srv.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpect")
		}
	}

	log.Println("Server exiting")
}
