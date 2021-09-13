package main

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"os/exec"
	"time"

	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
)

// Create struct to recive JSON format
type responseStruct struct {
	DeviceID uint32 `json:"device_id"`
}

/*
	Get local ip of this RPI
*/
func get_ip() net.IP {

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Logger.Info().Err(err)
		log.Logger.Error().Msgf("[X]\tConnection is down! [ERROR] - ", err)
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

/*
	Runs update command to update RPI
*/
func updateSystem() {
	log.Logger.Info().Msg("[+]\tFetching updates!")
	// fmt.Println("[+] Fetching updates!")
	cmd := exec.Command("bash", "-c", "sudo apt update && sudo apt upgrade -y && sudo apt autoremove -y &> /dev/null")
	// cmd.Stderr = os.Stdout
	// cmd.Stdout = os.Stdout
	err := cmd.Run()
	log.Logger.Info().Msgf("[+]\t[DONE] Updating")
	if err != nil {
		log.Logger.Error().Msgf("[X]\tCommand running failed [ERROR] - ", err)
	}
}

/*
	Checks if RPI has internet
*/
func checkForInternet() {

	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		log.Logger.Info().Err(err)
		log.Logger.Error().Msgf("[X]\tConnection is down!")
	} else {
		log.Logger.Info().Msgf("[+]\tConnection is up!")
		log.Logger.Info().Msgf("[*]\tIP is -> %s", get_ip())
		updateSystem()
		defer conn.Close()
	}
}

/*
	Returns URL
*/
func getURLForC2Server() string {

	c2_host := "127.0.0.1"
	url := "http://" + c2_host + ":8000/api/devices/addDevice"

	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", c2_host+":8000", timeout)
	if err != nil {
		log.Logger.Error().Msgf("[+]\tSite unreachable, [ERROR] -", err)
		log.Logger.Fatal()
	}
	log.Logger.Info().Msgf("[+]\tC2 is Alive -> %s", conn.LocalAddr().String())

	return url
}

/*
	Makes API call to C&C server
	Makes a post resquest to API
	Receives JSON data with DeviceID for RPI
*/
func api_call_addDevice() uint32 {

	ipAddr := get_ip().String()
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + "XxPFUhQ8R7kKhpgubt7v"

	//Encode the data
	postBody, _ := json.Marshal(map[string]string{
		"ip_str": ipAddr,
	})
	responseBody := bytes.NewBuffer(postBody)

	// Create a new request using http
	req, err := http.NewRequest("POST", getURLForC2Server(), responseBody)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError on response.\n[ERROR] -", err)

	}
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError on response.\n[ERROR] -", err)

	}

	var respStruct responseStruct

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&respStruct); err != nil {
		log.Logger.Error().Msgf("[X]\tError in decode.\n[ERROR] -", err)
	}
	log.Logger.Info().Msgf("[*]\tNew DeviceID -> %d", respStruct.DeviceID)
	defer resp.Body.Close()

	log.Logger.Info().Msgf("[+]\tDONE")

	return respStruct.DeviceID
}

/*
	Runs fucntions in order.
*/
func main() {
	log.InitLog(true)
	// checkForInternet()
	createConfigFile()
	// fmt.Println("\n [+] Server running!")
	// http.HandleFunc("/", handler)
	// log.Fatal(http.ListenAndServe(":8080", nil))

}
