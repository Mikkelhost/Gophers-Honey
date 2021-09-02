package httpserver

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "login.html", nil)
	} else if r.Method == "POST"{
		err := r.ParseForm()
		if err != nil {
			w.Write([]byte("Error in html form"))
		}
		username := r.FormValue("username")
		log.Debug().Msgf("Username: %s", username)
		password := r.FormValue("password")
		log.Debug().Msgf("Password: %s", password)
		w.Write([]byte("Handling userinformation"))
	} else {
		w.Write([]byte("Unsupported method"))
	}
}

func registrationHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "register.html", nil)
	}
}
