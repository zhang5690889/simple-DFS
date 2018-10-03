package main

import (
	"encoding/json"
	"log"
	"strconv"
	"net/http"
	"bytes"
	"io/ioutil"
	"fmt"
)

type client struct {
	proxy_port int
}

// user request data of the given key
func (c client) client_get(key string) {
	// steps
	// 1. client sends key to proxy
	// 2. proxy sends info of the servers which contain the key
	// 3. client sends request
	// 4. server sends data back

	// Step1  send meta data to proxy
	log.Printf("[Client]Sending key:%s to proxy.", key)
	jsonValue, _ := json.Marshal(map[string]string{"key": key})
	url := "http://127.0.0.1:" + strconv.Itoa(c.proxy_port) + "/get"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		log.Fatal(err)
	}

	//Step 2
	body, _ := ioutil.ReadAll(resp.Body)
	var server_urls []string;
	json.Unmarshal(body, &server_urls)
	log.Printf("[Client]Received server ports. Requesting key:%s", key)

	//Step 3 and 4
	retrieved_value := ""
	for _, server_url := range server_urls {

		// sends the key to servers
		resp, err = http.Post(server_url, "application/json", bytes.NewBuffer(jsonValue))

		if err != nil {
			log.Fatal(err)
		}
		body, _ = ioutil.ReadAll(resp.Body)
		var data KeyValueContainer;
		json.Unmarshal(body, &data)
		retrieved_value += data.Value
	}

	log.Printf("[Client]Client successfully retrieved data for key: %s value: %s", key, retrieved_value)
}

// user wants to save key value pair. (NO duplicate)
func (c client) client_put(data KeyValueContainer) {

	// steps
	// 1. client sends meta data to proxy
	// 2. proxy sends Allocation scheme back
	// 3. client cuts files into pieces
	// 4. client sends out data chunks to each server based on Allocation scheme

	log.Printf("[Client]Received value:%s. Sending to proxy for allocation scheme", data.Value)

	// Step1  send meta data to proxy
	jsonValue, _ := json.Marshal(map[string]string{"Size": strconv.Itoa(len(data.Value))})
	url := "http://127.0.0.1:" + strconv.Itoa(c.proxy_port) + "/set"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		log.Fatal(err)
	}

	// Step 2
	body, _ := ioutil.ReadAll(resp.Body)
	var allocation_scheme []AllocationBlock;
	json.Unmarshal(body, &allocation_scheme)

	//step 3
	//pre processing data
	// this var stores sliced string
	data_chunk_result := []string{}
	current_cursor := 0
	for i := 0; i < len(allocation_scheme); i++ {
		current_data_chunk_size, err := strconv.Atoi(allocation_scheme[i].Size)
		if err != nil {
			// handle error
			fmt.Println(err)
		}
		data_chunk := data.Value[current_cursor : current_cursor+current_data_chunk_size]
		data_chunk_result = append(data_chunk_result, data_chunk)
		current_cursor += current_data_chunk_size
	}

	log.Printf("[Client]Sliced data into: %s.Sending to server...", data_chunk_result)
	// step 4
	for i := 0; i < len(data_chunk_result); i++ {
		to_server_data := KeyValueContainer{data.Key, data_chunk_result[i]}
		jsonValue, _ = json.Marshal(to_server_data)
		request_url := allocation_scheme[i].Server_url + "/server_set"
		resp, err = http.Post(request_url, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("[Client] Save operation succesful!")
}
