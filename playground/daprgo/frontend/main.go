package main

import (
	"bytes"
	// "context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
)

// The state store
const dapr_store_name = "statestore"

// A function that decodes a json object to a struct
func getJson(r io.Reader, target interface{}) error {
	return json.NewDecoder(r).Decode(target)
}

type Reservation struct {
	Username string
	Userid   string
	Hotelid  string
}

func getset(w http.ResponseWriter, r *http.Request) {

	// Get DAPR environment variables to setup URL
	var DAPR_HOST, DAPR_HTTP_PORT string
	var okHost, okPort bool
	if DAPR_HOST, okHost = os.LookupEnv("DAPR_HOST"); !okHost {
		DAPR_HOST = "http://localhost"
	}
	if DAPR_HTTP_PORT, okPort = os.LookupEnv("DAPR_HTTP_PORT"); !okPort {
		DAPR_HTTP_PORT = "3501"
	}

	// TODO: How to check that body is application/json
	var reservation Reservation
	err := getJson(r.Body, &reservation)
	if err != nil {
		panic(err)
	}

	var mucacheOn string

	mucacheOn = r.Header.Get("mucache")

	// TODO: Understand what Context does precisely
	// ctx := context.Background()

	// Encode the reservation to send to the backend service
	var send_buffer bytes.Buffer
	err = json.NewEncoder(&send_buffer).Encode(reservation)
	if err != nil {
		log.Fatalln(err)
	}

	// Perform the request
	client := &http.Client{}
	req, err := http.NewRequest("POST", DAPR_HOST+":"+DAPR_HTTP_PORT+"/getset", &send_buffer)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// Adding app-id as part of the header
	req.Header.Add("dapr-app-id", "goapp")

	// Debugging info for mucache header
	if (mucacheOn) == "on" {
		var parentRID = r.Header.Get("RID")
		req.Header.Add("parentRID", parentRID)
		log.Println("Mucache header pass through appHttp")
	}
	// Adding the mucache header
	req.Header.Add("mucache", "on")

	// Invoking a service
	response, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// Result handling
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Order passed: ", string(result))

	w.Write(result)

}

func main() {
	fmt.Println(runtime.GOMAXPROCS(1))
	http.HandleFunc("/getset", getset)
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		panic(err)
	}
}
