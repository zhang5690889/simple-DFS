package main

import (
	"net/http"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"math/rand"
	"reflect"
	"time"
	"fmt"
	"log"
)

// This method returns info of server which contains the requested key
func proxy_get(w http.ResponseWriter, r *http.Request) {

	//steps:
	// 1. Broadcast requests key value
	// 2. servers that contains the key will respond
	// 3. proxy sends the server port back

	// sends all server ports back
	all_server_urls := []string{}
	for _, port := range _server_ports {
		all_server_urls = append(all_server_urls, "http://127.0.0.1:"+strconv.Itoa(port)+"/server_get")
	}
	jsonValue, _ := json.Marshal(all_server_urls)
	w.Write(jsonValue)
	log.Printf("[Proxy]Received key. Sending server ports")
}

//shuffle an array
func Shuffle(slice interface{}) {
	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	length := rv.Len()
	for i := length - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		swap(i, j)
	}
}

// Alg for uniform distribution
func uniformDistribution(inputsize int, chunk int) []int {
	result := []int{}

	// step 1: distribute until reminder Etc, 12/5 each node gets 2
	// step 2: fill reminder
	// step 3: shuffle result and filled reminder
	// step 4: add them up to ensure randomness
	equal_size := inputsize / chunk
	for i := 0; i < chunk; i++ {
		result = append(result, equal_size)
	}

	reminder := inputsize % chunk
	filled_reminder := []int{}
	for i := 0; i < reminder; i++ {
		filled_reminder = append(filled_reminder, 1)
	}
	for i := 0; i < chunk-reminder; i++ {
		filled_reminder = append(filled_reminder, 0)
	}

	//shuffle both result and add them up
	Shuffle(result)
	Shuffle(filled_reminder)
	for i := 0; i < chunk; i++ {
		result[i] += filled_reminder[i]
	}
	return result
}

// create allocation scheme
func create_allocation_scheme(input_size int, server_ports []int) []AllocationBlock {
	allocation_scheme_result := []AllocationBlock{}
	distribution_size := uniformDistribution(input_size, len(server_ports))

	//	create block
	for i := 0; i < len(server_ports); i++ {
		allocationblock := AllocationBlock{
			"http://127.0.0.1:" + strconv.Itoa(server_ports[i]),
			strconv.Itoa(distribution_size[i]),
		}
		allocation_scheme_result = append(allocation_scheme_result, allocationblock)
	}
	return allocation_scheme_result
}

func proxy_set(w http.ResponseWriter, r *http.Request) {

	// receives meta data
	body, _ := ioutil.ReadAll(r.Body)
	var meta_data map[string]string;
	json.Unmarshal(body, &meta_data)

	input_data_size, err := strconv.Atoi(meta_data["Size"])
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	// calculate allocation scheme
	allocation_scheme := create_allocation_scheme(input_data_size, _server_ports)

	// send Allocation scheme to client
	jsonValue, _ := json.Marshal(allocation_scheme)
	w.Write(jsonValue)
	log.Printf("[Proxy]Received data size:%d. Created scheme: %s", input_data_size, allocation_scheme)
}

// Proxy server should know how to connect to all servers at all time
var _server_ports []int

func start_proxy(proxy_port int, server_ports []int) {
	rand.Seed(time.Now().Unix())
	//save server ports in proxy's memory
	_server_ports = server_ports

	proxy_service := http.NewServeMux()
	// proxy only process meta data.
	proxy_service.HandleFunc("/get", proxy_get)
	proxy_service.HandleFunc("/set", proxy_set)
	http.ListenAndServe(":"+strconv.Itoa(proxy_port), proxy_service)
}
