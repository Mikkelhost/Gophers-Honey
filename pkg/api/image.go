package api

import (
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
)

func imageSubRouter(r *mux.Router) {
	imageRouter := r.PathPrefix("/images").Subrouter()
	//imageRouter.Handle("/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	imageRouter.Queries("download","{download:[0-9]+}").HandlerFunc(downloadImage).Methods("GET", "OPTIONS").Name("download")

}

func downloadImage(w http.ResponseWriter, r *http.Request) {
	image := mux.Vars(r)["download"]
	log.Logger.Debug().Msgf("User wants to download image: %s", image)
	w.Header().Set("Content-Disposition", "inline; filename=raspberrypi.img")
	http.ServeFile(w, r, "./images/"+image+".img")
}
