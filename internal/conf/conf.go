package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/conf"
)

// WebServer holds the server configuration settings
type WebServer struct {
	Host              string        `conf:"default:localhost:8080"`
	ReadTimeout       time.Duration `conf:"default:5s"`
	WriteTimeout      time.Duration `conf:"default:30s"`
	IdleTimeout       time.Duration `conf:"default:30s"`
	ReadHeaderTimeout time.Duration `conf:"default:5s"`
	ShutdownTimeout   time.Duration `conf:"default:5s"`
}

// Logging hold the log configuration settings
type Logging struct {
	Prefix string `conf:"default:[WG] :"`
}

// SMTP hold the smtp connection configuration settings
type SMTP struct {
	Server   string `conf:"default:smtp.sendgrid.net"`
	Port     int    `conf:"default:587"`
	From     string `conf:"default:WG Concierge<no-reply@wg-concierge.com>"`
	Username string `conf:"required"`
	Password string `conf:"required,noprint" json:"-"`
}

// App holds the application configuration
type App struct {
	Web    WebServer
	Log    Logging
	Email  SMTP
	Domain string `conf:"default:http://localhost:8080"`
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

func PrettyString(app *App) string {
	b, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(b)
}
