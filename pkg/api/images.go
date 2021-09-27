package api

import (
	"encoding/json"
	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"github.com/gorilla/mux"
	"net/http"
)

func imageSubRouter(r *mux.Router) {
	imageAPI := r.PathPrefix("/api/images").Subrouter()
	//imageRouter.Handle("/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	imageAPI.Queries("download","{download:[0-9]+}").HandlerFunc(tokenAuthMiddleware(downloadImage)).Methods("GET", "OPTIONS").Name("download")
	imageAPI.HandleFunc("/getImages", tokenAuthMiddleware(getImages))

}

func downloadImage(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	image := mux.Vars(r)["download"]
	log.Logger.Debug().Msgf("APIUser wants to download image: %s", image)
	w.Header().Set("Content-Disposition", "inline; filename=raspberrypi.img")
	http.ServeFile(w, r, "./images/"+image+".img")
}

func getImages(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	var images []model.Image
	images, err := database.GetImages()
	if err != nil {
		w.Write([]byte("Error retrieving devices"))
		return
	}
	if len(images) == 0 {
		w.Write([]byte("No devices in DB"))
		return
	}
	json.NewEncoder(w).Encode(images)
}