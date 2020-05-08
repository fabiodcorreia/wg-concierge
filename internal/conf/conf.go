package conf

import (
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/conf"
)

// WebServer holds the server configuration settings
type WebServer struct {
	Host            string        `conf:"default:127.0.0.1:8080"`
	RequestTimeout  time.Duration `conf:"default:60s"`
	ReadTimeout     time.Duration `conf:"default:5s"`
	WriteTimeout    time.Duration `conf:"default:5s"`
	ShutdownTimeout time.Duration `conf:"default:5s"`
}

// Logging hold the log configuration settings
type Logging struct {
	Prefix string `conf:"default:WG:"`
}

// SMTP hold the smtp connection configuration settings
type SMTP struct {
	Server   string `conf:"default:smtp.sendgrid.net"`
	Port     int    `conf:"default:465"`
	From     string `conf:"defautl:WG Concierge<no-reply@wg-concierge.com>"`
	Username string `conf:"required"`
	Password string `conf:"required,noprint"`
}

// App holds the application configuration
type App struct {
	Web   WebServer
	Log   Logging
	Email SMTP
}

// Load will load the configuration into the App Configuration
func Load(args []string, namespace string, appConfig *App) error {
	if err := conf.Parse(args[1:], namespace, appConfig); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage(namespace, appConfig)
			if err != nil {
				return fmt.Errorf("generating config usage: %w", err)
			}
			fmt.Println(usage)
			os.Exit(0)
		}
		return fmt.Errorf("parsing config: %w", err)
	}
	return nil
}
