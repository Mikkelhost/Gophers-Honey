package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
)

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
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println("\n [+] Server running!")
	checkForInternet()
}
