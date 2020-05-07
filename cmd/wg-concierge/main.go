package main

import (

	// Register the expvar handlers

	"log"
	"net/http"
	"os"

	"github.com/fabiodcorreia/wg-concierge/internal/conf"
	"github.com/fabiodcorreia/wg-concierge/internal/web"
)

var version = "development"

func main() {
	if err := run(); err != nil {
		log.Println("error :", err)
		os.Exit(1)
	}
}

func run() error {

	// Load Configuration
	var cfg conf.App
	err := conf.Load(os.Args, "WG", &cfg)
	if err != nil {
		return err
	}

	// Setup Logger
	log := log.New(os.Stdout, cfg.Log.Prefix, log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)

	// Setup Service
	service := web.NewService(log)
	service.AddTimeout(cfg.Web.RequestTimeout)

	service.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Println("hello")
		w.Write([]byte("hello"))
	})

	// Call the NewAPI that receives a logger, external connections, like email, db, auth... and return an http.handler
	// Inside Call the NewApp that receives logger, this will return a chi.NewMux with the default middlewares logger, database, ....
	// Each domain will have a struct with specific resources like db and auth this struct will have a method for each service operation GET List, PUT,...
	// The App allows to add routes from teh domain and specific middlewares for the routes

	// Setup server
	svr := web.NewHTTPServer(cfg, log, service)

	return svr.StartAndWait(cfg.Web.ShutdownTimeout)
}
