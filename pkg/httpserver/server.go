package httpserver

import (
	"fmt"
	"github.com/GeekMuch/GoHoney/pkg/database"
	"html/template"
	"log"
	"net/http"
	"os"
	logger "github.com/rs/zerolog/log"
	"github.com/GeekMuch/GoHoney/pkg/websocket"
	"github.com/gorilla/mux"
)

var tpl *template.Template
var configured bool
var wd = GetWd()

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.HandleFunc("/favicon.ico", faviconHandler)
	r.HandleFunc("/setup", setupHandler)
	r.HandleFunc("/", homePage)
	r.HandleFunc("/devices", devicesHandler)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/register", registrationHandler)
	r.PathPrefix("/static/").Handler(s)

	return r
}
func faviconHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "/static/img/favicon.ico")

}

func initTpl() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func GetWd() string{
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return wd
}
func parseTemplates(givenTemplate string) (*template.Template, error) {
	var tmpl *template.Template
	var err error
	tmpl, err = template.ParseFiles(
		wd + "/templates/base.html",
		wd + "/templates/navbar.html",
		wd + givenTemplate,
	)


	return tmpl, err
}

func RunServer() {
	fmt.Println("Starting websocket")
	router := setupRouter()
	websocket.SetupRouter(router)
	initTpl()
	log.Fatal(http.ListenAndServeTLS(":8443", "certs/nginx-selfsigned.crt", "certs/nginx-selfsigned.key", router))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	config := database.GetConf() // Has the service setup been run
	if config {
		devicesHandler(w,r)
	} else {
		firstRun(w, r)
	}
}

func firstRun(w http.ResponseWriter, r *http.Request) {
	pd.Navbar = false
	tmpl, err := parseTemplates("/templates/firstrun.html")
	if err != nil {
		logger.Error().Msgf("Error: %s", err)
	}
	tmpl.Execute(w, pd)
}

func setupHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" {
		r.ParseForm()
		for key, values := range r.Form{
			for _, value := range values {
				log.Println(key, value)
			}
		}


	} else {
		w.Write([]byte("Invalid http method"))
	}
}



func devicesHandler(w http.ResponseWriter, r *http.Request) {
	pd.Navbar = true
	tmpl, err := parseTemplates("/templates/devices.html")
	if err != nil {
		log.Println("error index tmpl: ", err)
	}
	tmpl.Execute(w, pd)
}