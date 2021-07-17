package service

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/fabiodcorreia/wg-concierge/internal/qr"
	"github.com/fabiodcorreia/wg-concierge/internal/repository"
	"github.com/fabiodcorreia/wg-concierge/internal/wg"
	"github.com/fabiodcorreia/wg-concierge/internal/wg/serializer"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type RegistrationService struct {
	*chi.Mux
	logger *log.Logger
	repo   *repository.Repository
	//templates Templates
}

func NewRegistrationService(logger *log.Logger, repo *repository.Repository) http.Handler {
	s := RegistrationService{
		Mux:    chi.NewRouter(),
		logger: logger,
		repo:   repo,
	}
	s.Get("/{token}", s.handleRegistration())
	s.Post("/{token}", s.handleRegistration())
	return s
}

func (s *RegistrationService) handleRegistration() http.HandlerFunc {
	var init sync.Once
	var tpl *template.Template
	var tplerr error
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			s.logger.Println("registration-form compiled")
			tpl, tplerr = template.New("registration-form").Parse(templateRegistrationForm)
		})
		if tplerr != nil {
			http.Error(w, tplerr.Error(), http.StatusInternalServerError)
			return
		}
		key := chi.URLParam(r, "token")
		token, err := uuid.Parse(key)
		if err != nil {
			//TODO log and return error
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		i, err := s.repo.GetInvitation(token)
		if err != nil {
			fmt.Println(err)
			//TODO log and return error
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// This block should be on the Post not here on the Get
		cc := wg.ClientConfig{
			Interface: wg.ClientInterface{
				PrivateKey: "", // Generate Private Key
				Address:    "", // Get the next available Address
			},
			Peer: wg.ClientPeer{
				Endpoint:   "vpnserver",      // From Configuration
				PublicKey:  "server pub key", // From Server Config File
				AllowedIPs: "0.0.0.0",        // All by default
			},
		}

		// QR should have a separated URL that return the image to the page
		cs, err := serializer.MarshalClientToStr(cc)
		qr.EncodeString(cs, w)

		fmt.Println(i)

		// TODO s.repo.RevokeInvitation(token)
		tpl.Execute(w, nil)
	}
}

const templateRegistrationForm = `
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
					<label for="InputDevice">Device Name</label>
					<input type="text" class="form-control" id="device" name="device" aria-describedby="deviceHelp" required>
					</div>
					<button type="submit" class="btn btn-primary">
						<i class="fas fa-save"></i> Invite
					</button>
				</form>
			</div>
		</div>
	</div>
  </body>
</html>
`

const templateRegistrationFail = `
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
