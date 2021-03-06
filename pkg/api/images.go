package api

import (
	"encoding/json"
	"fmt"
	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"github.com/Mikkelhost/Gophers-Honey/pkg/piimage"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

//imageSubRouter
//Routes image related api endpoints to their respective handlers
func imageSubRouter(r *mux.Router) {
	imageAPI := r.PathPrefix("/api/images").Subrouter()
	//imageRouter.Handle("/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	imageAPI.Queries("download", "{download:[0-9]+}").HandlerFunc(tokenAuthMiddleware(downloadImage)).Methods("GET", "OPTIONS").Name("download")
	imageAPI.HandleFunc("", tokenAuthMiddleware(imageHandler)).Methods("GET", "POST", "DELETE", "OPTIONS")

}

//imageHandler
//Checks the request method to determine which handler to use.
func imageHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	// CORS preflight handling.
	if r.Method == "OPTIONS" {
		return
	}
	switch r.Method {
	case "GET":
		getImages(w, r)
		return
	case "POST":
		newImage(w, r)
		return
	case "DELETE":
		removeImage(w, r)
		return
	}
}

//downloadImage
//Serves the requested raspberry pi image to a user for download
func downloadImage(w http.ResponseWriter, r *http.Request) {
	image := mux.Vars(r)["download"]
	log.Logger.Debug().Msgf("APIUser wants to download image: %s", image)
	w.Header().Set("Content-Disposition", "inline; filename=raspberrypi.img")
	http.ServeFile(w, r, "./images/"+image+".img")
}

//getImages
//Gets all current images
func getImages(w http.ResponseWriter, r *http.Request) {
	var images []model.Image
	images, err := database.GetImages()
	if err != nil {
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Error retrieving images"})
		return
	}
	if len(images) == 0 {
		json.NewEncoder(w).Encode(model.APIResponse{Error: "No images in DB"})
		return
	}
	json.NewEncoder(w).Encode(images)
}

//newImage
//Creates a new raspberry pi image which later can be downloaded through the images/download endpoint
func newImage(w http.ResponseWriter, r *http.Request) {
	var imgInfo = model.ImageInfo{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&imgInfo); err != nil {
		log.Logger.Warn().Msgf("Failed decoding json: %s", err)
		w.Write([]byte(fmt.Sprintf("Failed decoding json: %s", err)))
		return
	}
	log.Logger.Debug().Msgf("setup params: %v", imgInfo)

	port, err := strconv.Atoi(imgInfo.Port)
	if err != nil {
		log.Logger.Warn().Msgf("Error converting port to int: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("%s", err)})
		return
	}
	if port < 0 {
		log.Logger.Warn().Msg("Port smaller than 0, aboritng")
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Port cannot be smaller than 0"})
		return
	}

	id, err := database.NewImage(model.Image{
		Name: imgInfo.ImageName,
	})
	if err != nil {
		log.Logger.Warn().Msgf("Error creating new image in db: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("%s", err)})
		return
	}

	err = piimage.InsertConfig(model.PiConf{
		C2:        imgInfo.C2,
		Port:      port,
		DeviceID:  0,
		DeviceKey: DEVICE_KEY,
	}, id)
	if err != nil {
		log.Logger.Warn().Msgf("Error inserting config into image: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("%s", err)})
		return
	}
	json.NewEncoder(w).Encode(model.APIResponse{Error: ""})
}

//removeImage
//Removes a image both from disk and database.
func removeImage(w http.ResponseWriter, r *http.Request) {
	var image model.Image
	var err error

	decoder := json.NewDecoder(r.Body)

	if err = decoder.Decode(&image); err != nil {
		log.Logger.Warn().Msgf("Error decoding JSON: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding JSON: %s", err)})
		return
	}
	// Remove image from DB collection.
	log.Logger.Debug().Uint32("id", image.Id).Msg("Deleting image from db")
	if err = database.RemoveImage(image.Id); err != nil {
		log.Logger.Warn().Msgf("Error removing image from collection: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error removing image from db: %s", err)})
		return
	}
	// Delete image file from disk
	log.Logger.Debug().Uint32("id", image.Id).Msg("Deleting image from disk")
	if err = piimage.DeleteImage(image.Id); err != nil {
		log.Logger.Warn().Msgf("Error deleting image from disk: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error deleting image from disk: %s", err)})
		return
	}

	json.NewEncoder(w).Encode(model.APIResponse{Error: ""})
}
