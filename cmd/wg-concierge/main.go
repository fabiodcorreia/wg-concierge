package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fabiodcorreia/wg-concierge/internal/conf"
	"github.com/fabiodcorreia/wg-concierge/internal/email"
	"github.com/fabiodcorreia/wg-concierge/internal/repository"
	"github.com/fabiodcorreia/wg-concierge/internal/web"
	"github.com/fabiodcorreia/wg-concierge/internal/web/service"
)

var version = "development"

func main() {
	if err := run(); err != nil {
		log.Fatalln("error :", err)
	}
}

func run() error {

	// Load Configuration
	var cfg conf.App
	cfg.Version.SVN = version
	cfg.Version.Desc = `WG-Concierge`
	err := conf.Load(os.Args, "WG", &cfg)
	if err != nil {
		return err
	}

	// Setup Logger
	log := log.New(os.Stdout, cfg.Log.Prefix, log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix)

	//TODO check if the wireguard configuration file can be found and read

	// Show Version
	log.Printf("Started : Initializing : version %q", cfg.Version.SVN)

	// Setup Email Sender
	emailSender := email.NewSender(cfg.Email.Username, cfg.Email.Password, cfg.Email.Server, cfg.Email.From, cfg.Email.Port)

	// Setup Database
	if cfg.Database == "$HOME/.wg-concierge" {
		path, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("user home directory not found: %w", err)
		}
		cfg.Database = fmt.Sprintf("%s/.wg-concierge", path)
	}

	repo, err := repository.NewRepository(cfg.Database)
	if err != nil {
		return fmt.Errorf("fail to stablish database connection: %w", err)
	}
	defer repo.Close()

	// Show Application Configuration
	log.Println(cfg.PrettyString())

	// Setup server
	svr := web.NewHTTPServer(cfg, log)
	svr.AddService("/invite", service.NewInvitationService(log, cfg, emailSender, repo))
	svr.AddService("/register", service.NewRegistrationService(log, repo))

	return svr.StartAndWait(cfg.Web.ShutdownTimeout)
}
