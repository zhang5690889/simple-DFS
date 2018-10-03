package main

import (
	"net/http"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"log"
)

// This variable hold in memory storage
// It seems that go doesn't support thread global variable
// Have to get around it using a hack. Use dictionary to store data for each http service
var in_memory_storage map[string][]KeyValueContainer

// This method saves the data segments from client
func server_set(w http.ResponseWriter, r *http.Request) {

	// receive data
	body, _ := ioutil.ReadAll(r.Body)
	var data KeyValueContainer
	json.Unmarshal(body, &data)
	log.Printf("[Server] %s : Recived data:%s (size:%d)\n", r.Host, data.Value, len(data.Value))

	_, ok := in_memory_storage[r.Host]
	if ok {
		in_memory_storage[r.Host] = append(in_memory_storage[r.Host], data)

	} else {
		in_memory_storage[r.Host] = []KeyValueContainer{data}
	}
}

// This method returns data found in the in memory storage
func server_get(w http.ResponseWriter, r *http.Request) {

	// receive key
	body, _ := ioutil.ReadAll(r.Body)
	var data map[string]string
	json.Unmarshal(body, &data)

	//search for key and return
	server_data := in_memory_storage[r.Host]
	for index := range server_data {
		if server_data[index].Key == data["key"] {
			jsonValue, _ := json.Marshal(server_data[index])
			w.Write(jsonValue)
			log.Printf("[Server]Received request for key: %s. Sending value : %s", data["key"], server_data[index].Value)
			break
		}
	}
}

// This method starts server services
func start_server(server_port int) {
	// initialize in memory database
	in_memory_storage = make(map[string][]KeyValueContainer)

	// new service. GO thing
	server_service := http.NewServeMux()
	server_service.HandleFunc("/server_get", server_get)
	server_service.HandleFunc("/server_set", server_set)
	http.ListenAndServe(":"+strconv.Itoa(server_port), server_service)
}
