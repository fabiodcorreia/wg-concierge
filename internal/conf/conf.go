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
	From     string `conf:"default:WG Concierge<no-reply@wg-concierge>"`
	Username string `conf:"required"`
	Password string `conf:"required,noprint" json:"-"`
}

type WireGuard struct {
	ConfigPath     string `conf:"default:/etc/wireguard/wg0.conf"`
	DNS            string `conf:"default:1.1.1.1"`
	PublicEndpoint string `conf:"required"`
}

// App holds the application configuration
type App struct {
	conf.Version
	Web       WebServer
	Email     SMTP
	AdminUser string `conf:"default:admin"`
	AdminPass string `conf:"default:admin,noprint" json:"-"`
	Domain    string `conf:"default:http://localhost:8080"`
	Database  string `conf:"default:$HOME/.wg-concierge"`
	WireGuard WireGuard
	Log       Logging
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
		if err == conf.ErrVersionWanted {
			version, err := conf.VersionString(namespace, appConfig)
			if err != nil {
				return fmt.Errorf("generating version: %w", err)
			}
			fmt.Println(version)
			os.Exit(0)
		}
		return fmt.Errorf("parsing config: %w", err)
	}
	return nil
}

// PrettyString generates a JSON well formated string representation of the configuraiton
func (app *App) PrettyString() string {
	b, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(b)
}
