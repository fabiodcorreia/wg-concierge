package main

import (
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/fabiodcorreia/wg-concierge/internal/conf"
	"github.com/fabiodcorreia/wg-concierge/internal/email"
	"github.com/fabiodcorreia/wg-concierge/internal/web"
	"github.com/go-chi/chi"
)

var version = "development"

func main() {
	if err := run(); err != nil {
		log.Println("error :", err)
		os.Exit(1)
	}
}

var emailRegex = regexp.MustCompile(`[^@]+@[^\.]+\..+`)

func run() error {

	// Load Configuration
	var cfg conf.App
	err := conf.Load(os.Args, "WG", &cfg)
	if err != nil {
		return err
	}

	// Setup Logger
	log := log.New(os.Stdout, cfg.Log.Prefix, log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)

	// Setup Email Sender

	emailSender := email.NewSender(cfg.Email.Username, cfg.Email.Password, cfg.Email.Server, cfg.Email.From, cfg.Email.Port)

	// Setup Service
	service := web.NewService(log)
	service.AddTimeout(cfg.Web.RequestTimeout)
	service.Get("/invite", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	service.Post("/invite", func(w http.ResponseWriter, r *http.Request) {
		//to, found := r.Context().Value(middleware.URLFormatCtxKey).(string)
		to := chi.URLParam(r, "email")

		if to == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Not Found")
			w.Write([]byte("Email not found"))
			return
		}
		log.Println(emailRegex.Match([]byte(to)))
		if emailRegex.Match([]byte(to)) {
			log.Println(emailSender.SendInvitation(to, ""))
			w.Write([]byte("Sent"))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Invalid Email")
			log.Println(to)
			w.Write([]byte("Invalid Email found"))
		}

	})

	// Setup server
	svr := web.NewHTTPServer(cfg, log, service)

	return svr.StartAndWait(cfg.Web.ShutdownTimeout)
}
