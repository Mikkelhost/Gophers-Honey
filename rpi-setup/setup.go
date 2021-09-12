package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func get_ip() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Print(err)
		fmt.Println("[X] Connection is down!")
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func runCommand() {
	cmd := exec.Command("bash", "-c", "sudo apt update && sudo apt upgrade -y")
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}
func updateSystem() {
	fmt.Println("[+] Fetching updates!")
	runCommand()
}

// Get preferred outbound ip of this machine
func checkForInternet() {

	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		log.Print(err)
		fmt.Println("[X] Connection is down!")
	} else {
		fmt.Println("[+] Connection is up!")
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		fmt.Printf("[**] IP is -> %s\n", localAddr.IP)
		updateSystem()
		defer conn.Close()
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("ls", "-al")
	fmt.Fprintf(w, "%s %s\n", cmd.Run(), r.URL.Path[1:])
	fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL.Path[1:])

}

func api_call_addDevice() {

	ipAddr := get_ip().String()
	c2_host := "127.0.0.1"
	url := "http://" + c2_host + ":8000/api/devices/addDevice"

	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", c2_host+":8000", timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		log.Fatal(1)
	}
	fmt.Printf("\n [+] C2 is Alive -> %s", conn.LocalAddr().String())
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + "XxPFUhQ8R7kKhpgubt7v"

	//Encode the data
	postBody, _ := json.Marshal(map[string]string{
		"deviceKey": "XxPFUhQ8R7kKhpgubt7v",
		"ip_str":    ipAddr,
	})
	responseBody := bytes.NewBuffer(postBody)

	// Create a new request using http
	req, err := http.NewRequest("POST", url, responseBody)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)

	}
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)

	}
	type responseStruct struct {
		DeviceID uint32 `json:"device_id"`
	}
	var respStruct responseStruct

	decoder := json.NewDecoder(resp.Body)
	// decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&respStruct); err != nil {
		log.Println("Error in decode")
	}
	fmt.Printf("\n\n Respons -> %d", respStruct.DeviceID)
	fmt.Printf("\n\n NEW Respons -> %s", decoder)
	defer resp.Body.Close()

	// resp, err := http.Post(url,"application/json", responseBody)

	// //Handle Error
	// if err != nil {
	// 	log.Fatalf("An Error Occured %v", err)
	// }
	// defer resp.Body.Close()

	// //Read the response body
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// sb := string(body)
	// log.Printf(sb)
	// }
	fmt.Println("[+] DONE")
}

func main() {
	api_call_addDevice()
	// fmt.Println("\n [+] Server running!")
	// http.HandleFunc("/", handler)
	// log.Fatal(http.ListenAndServe(":8080", nil))
	// checkForInternet()
}
