package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/fabiodcorreia/wg-concierge/internal/conf"
	"github.com/fabiodcorreia/wg-concierge/internal/email"
	"github.com/fabiodcorreia/wg-concierge/internal/web"
	"github.com/google/uuid"
)

var version = "development"

func main() {
	if err := run(); err != nil {
		log.Fatalln("error :", err)
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

	log.Printf("Started : Initializing : version %q", version)
	log.Println(conf.PrettyString(&cfg))

	// Setup Email Sender

	emailSender := email.NewSender(cfg.Email.Username, cfg.Email.Password, cfg.Email.Server, cfg.Email.From, cfg.Email.Port)

	// Setup Dummy Repository
	repository := make(map[string]string)

	// Setup Service
	service := web.NewService(log)

	service.Get("/invite", handleInviteUser())

	service.Post("/invite", func(w http.ResponseWriter, r *http.Request) {

		// Generate a UUID to send on the Email
		// Store that UUID on the repository
		key, err := uuid.NewRandom()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			w.Write([]byte("Fail to generate invitation key"))
			return
		}

		err = r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("ErrOR")
			w.Write([]byte("ErrOR ErrOR ErrOR"))
			return
		}
		to := r.Form.Get("email")

		if to == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Not Found")
			w.Write([]byte("Email not found"))
			return
		}

		if emailRegex.Match([]byte(to)) {
			err := emailSender.SendInvitation(to, fmt.Sprintf("%s/register/%s", cfg.Domain, key))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			repository[key.String()] = to
			w.WriteHeader(http.StatusOK)
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

func handleInviteUser() http.HandlerFunc {
	var init sync.Once
	var tpl *template.Template
	var tplerr error
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tpl, tplerr = template.New("invite-form").Parse(templateInvite)
		})
		if tplerr != nil {
			http.Error(w, tplerr.Error(), http.StatusInternalServerError)
			return
		}
		tpl.Execute(w, nil)
	}
}

const templateInvite = `
<!doctype html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">

    <title>WG Concierge Invite</title>
  </head>
  <body>
	<div class="container">
		<div class="row">
			<div class="col">
			<form method="post">
			<div class="form-group">
			  <label for="InputEmail">Email Address</label>
			  <input type="email" class="form-control" id="email" name="email" aria-describedby="emailHelp" required>
			</div>
			<button type="submit" class="btn btn-primary">Invite</button>
		  </form>
			</div>
		</div>
	</div>
  </body>
</html>
`
