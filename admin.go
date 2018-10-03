package main

import (
	"log"
	"time"
)

// starting proxy service
func activate_proxy(proxy_port int, server_ports []int) {
	start_proxy(proxy_port, server_ports)
}

func activate_server(server_port int) {
	start_server(server_port)
}

// user command simulation
func simulate_user(proxy_port_num int) {

	//	create some sample data
	data1 := KeyValueContainer{"John", "abcdefghijklmn"}
	data2 := KeyValueContainer{"Mike", "some data about Mike"}

	//create client
	client := client{proxy_port_num}
	log.Printf("User load data. (key: %s value: %s )\n", data1.Key, data1.Value)
	client.client_put(data1)
	log.Printf("User load data. (key: %s value: %s )\n", data2.Key, data2.Value)
	client.client_put(data2)

	log.Printf("User retrieve key: %s \n", data1.Key)
	client.client_get(data1.Key)
	log.Printf("User retrieve key: %s \n", data2.Key)
	client.client_get(data2.Key)
}

//Test automation
func main() {

	// Configuration
	proxy_port := 1234
	server_ports := []int{1300, 1311, 1400}

	// Starts all services
	go activate_proxy(proxy_port, server_ports)

	for _, port := range server_ports {
		go activate_server(port)
	}

	log.Println("All Services started...")
	//wait 2 sends to make sure all services started
	time.Sleep(2 * time.Second)

	// User operations
	simulate_user(proxy_port)
}
