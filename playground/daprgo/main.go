package main

import (
	// "context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"

	// dapr "github.com/dapr/go-sdk/client"
)

// The state store
const dapr_store_name = "statestore"

func echo(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, HTTP!\n")
}

// A function that decodes a json object to a struct
func getJson(r io.Reader, target interface{}) error {
	return json.NewDecoder(r).Decode(target)
}

type Reservation struct {
	Username string
	Userid   string
	Hotelid  string
}

// {"user_name":"User Name","user_id":"2", "hotel_id":"10"}

func getset(w http.ResponseWriter, r *http.Request) {

	// Get DAPR environment variables to setup URL
	var DAPR_HOST, DAPR_HTTP_PORT string
	var okHost, okPort bool
	if DAPR_HOST, okHost = os.LookupEnv("DAPR_HOST"); !okHost {
		DAPR_HOST = "http://localhost"
	}
	if DAPR_HTTP_PORT, okPort = os.LookupEnv("DAPR_HTTP_PORT"); !okPort {
		DAPR_HTTP_PORT = "3500"
	}

	// TODO: How to check that body is application/json
	var reservation Reservation
	err := getJson(r.Body, &reservation)
	if err != nil {
		panic(err)
	}
	// log.Println(reservation)
	var mucacheOn bool

	_, mucacheOn = r.Header["Mucache"]

	// TODO: Understand what Context does precisely
	// ctx := context.Background()

	key := reservation.Hotelid
	client := &http.Client{}
	req, err := http.NewRequest("GET", DAPR_HOST+":"+DAPR_HTTP_PORT+"/v1.0/state/statestore/"+key, nil)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	req.Header.Add("Dapr-app-id", "statestore")

	if mucacheOn {
		req.Header.Add("mucache", "on")
		var parentRID = r.Header.Get("RID")
		req.Header.Add("parentRID", parentRID)
	}

	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Order passed: ", string(result))
	// log.Printf("data [key:%s etag:%s]: %s", item.Key, item.Etag, string(item.Value))

	// TODO: Save at statestore

	fmt.Fprintf(w, "Reservation from %s for hotel: %s succeeded.\n", reservation.Username, reservation.Hotelid)
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(1))
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/getset", getset)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
