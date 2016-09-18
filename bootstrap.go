package main

import (
	"net"
	"fmt"
	"encoding/json"
	"strconv"
	"bytes"
	"encoding/base64"
	"os"
	"net/http"
)

func bootstrapAndServe() {
	lsnr, err := net.Listen("tcp4", ":0")
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	fmt.Println("Listening on:", lsnr.Addr().String())

	go sendConfigToServer(lsnr.Addr())

	err = http.Serve(lsnr, nil)
	if err != nil {
		panic(err)
	}

}

func sendConfigToServer(addr net.Addr) {
	url := "http://" + globalSettings.BootstrapAddress + "/config/"
	fmt.Printf("Manager location: %s\n", url)

	jsonStr, _ := json.Marshal(ReplicatServer{Name: globalSettings.Name, Address: "127.0.0.1:" + strconv.Itoa(addr.(*net.TCPAddr).Port)})
	fmt.Printf("jsonStr: %s\n", jsonStr)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	data := []byte(globalSettings.ManagerCredentials)
	authHash := base64.StdEncoding.EncodeToString(data)
	req.Header.Add("Authorization", "Basic "+authHash)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&serverMap)
		fmt.Printf("configHander nodes: %s", serverMap)

		if err != nil {
			fmt.Println(err)
		}
	}
}