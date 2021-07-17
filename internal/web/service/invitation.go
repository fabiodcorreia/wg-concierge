package service

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/fabiodcorreia/wg-concierge/internal/conf"
	"github.com/fabiodcorreia/wg-concierge/internal/email"
	"github.com/fabiodcorreia/wg-concierge/internal/repository"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

type invitationTemplates struct {
	getInvite     *template.Template
	successInvite *template.Template
	failInvite    *template.Template
}

// InvitationService is responsible to handle all the invitation routes and actions
type InvitationService struct {
	*chi.Mux
	logger    *log.Logger
	repo      *repository.Repository
	templates invitationTemplates
}

// NewInvitationService creates a new service to handle the invitations with all the service dependencies
func NewInvitationService(logger *log.Logger, cfg conf.App, emailSender *email.Sender, repo *repository.Repository) http.Handler {
	s := InvitationService{
		Mux:    chi.NewRouter(),
		logger: logger,
		repo:   repo,
	}
	s.templates.getInvite = template.Must(template.New("invite-form").Parse(templateInvitationForm))
	s.templates.successInvite = template.Must(template.New("invite-success").Parse(templateInvitationSuccess))
	s.templates.failInvite = template.Must(template.New("invite-fail").Parse(templateInvitationFail))

	auth := make(map[string]string, 1)
	auth[cfg.AdminUser] = cfg.AdminPass

	s.Use(middleware.BasicAuth("WG-CG", auth))
	s.Get("/", s.handleIndex())
	s.Post("/", s.handleInviteUser(emailSender, cfg.Domain))
	return s
}

// handleIndex receives requests for GET /invite that shows the
func (s *InvitationService) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.templates.getInvite.Execute(w, nil)
	}
}

func (s *InvitationService) handleInviteUser(emailSender *email.Sender, domain string) http.HandlerFunc {
	emailRegex := regexp.MustCompile(`[^@]+@[^\.]+\..+`)

	return func(w http.ResponseWriter, r *http.Request) {
		key, err := uuid.NewRandom()
		if err != nil {
			s.logger.Printf("Fail to generate invitation key: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			s.templates.failInvite.Execute(w, nil)
			return
		}

		err = r.ParseForm()
		if err != nil {
			s.logger.Printf("Fail to parse invitation form: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			s.templates.failInvite.Execute(w, nil)
			return
		}

		to := r.Form.Get("email")
		if to == "" {
			s.logger.Printf("Fail to get email address for invitation, is empty\n")
			w.WriteHeader(http.StatusBadRequest)
			s.templates.failInvite.Execute(w, nil)
			return
		}

		if emailRegex.Match([]byte(to)) {
			err := s.repo.AddInvitation(repository.Invitation{
				Email:   to,
				Token:   key,
				Created: time.Now(),
			})
			if err != nil {
				s.logger.Printf("Fail to save invitation to database: %v\n", err)

				err = s.repo.RevokeInvitation(key)
				if err != nil {
					s.logger.Printf("Fail to rollback invitation from database: %v\n", err)
				}

				w.WriteHeader(http.StatusInternalServerError)
				s.templates.failInvite.Execute(w, nil)
				return
			}

			err = emailSender.SendInvitation(to, fmt.Sprintf("%s/register/%s", domain, key))
			if err != nil {
				s.logger.Printf("Fail to send invitation email: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				s.templates.failInvite.Execute(w, nil)
				return
			}

			s.templates.successInvite.Execute(w, to)
		} else {
			s.logger.Printf("Fail to handle email '%s', is not valid\n", to)
			w.WriteHeader(http.StatusBadRequest)
			s.templates.failInvite.Execute(w, nil)
		}
	}
}

const templateInvitationSuccess = `
<!doctype html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
	<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
	
	<script src="https://kit.fontawesome.com/a3b4cb126c.js" crossorigin="anonymous"></script>

	<style>
  		body {
    		padding-top: 60px;
  		}
  		@media (max-width: 980px) {
    		body {
      			padding-top: 1px;
    		}
  		}
	</style>

    <title>WG Concierge Invite</title>
  </head>
  <body>
	<div class="container">
		<div class="row">
			<div class="col">
				<h1>Thank You!</h1>
				<p>The invitation was sent to <b>{{ . }}</b></p>
				<a class="btn btn-primary" href="/invite" role="button">Home</a>
			</div>
		</div>
	</div>
  </body>
</html>
`

const templateInvitationFail = `
<!doctype html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
	<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
	
	<script src="https://kit.fontawesome.com/a3b4cb126c.js" crossorigin="anonymous"></script>

	<style>
  		body {
    		padding-top: 60px;
  		}
  		@media (max-width: 980px) {
    		body {
      			padding-top: 1px;
    		}
  		}
	</style>

    <title>WG Concierge Invite</title>
  </head>
  <body>
	<div class="container">
		<div class="row">
			<div class="col">
				<h1>Invitation Fail!</h1>
				<a class="btn btn-primary" href="/invite" role="button">Home</a>
			</div>
		</div>
	</div>
  </body>
</html>
`

const templateInvitationForm = `
<!doctype html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
	<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
	
	<script src="https://kit.fontawesome.com/a3b4cb126c.js" crossorigin="anonymous"></script>

	<style>
  		body {
    		padding-top: 60px;
  		}
  		@media (max-width: 980px) {
    		body {
      			padding-top: 1px;
    		}
  		}
	</style>

    <title>WG Concierge Invite</title>
  </head>
  <body>
	<div class="container">
		<div class="row">
			<div class="col">
			<form method="post">
			<div class="form-group">
			  <label for="InputEmail">Email Address</label>
			  <input type="email" class="form-control" id="email" name="email" aria-describedby="emailHelp" required placeholder="user@email.com">
			</div>
			<button type="submit" class="btn btn-primary">
				<i class="fas fa-paper-plane"></i> Invite
			</button>
		  </form>
			</div>
		</div>
	</div>
  </body>
</html>
`
