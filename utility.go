package main

// Type for storing key value
type KeyValueContainer struct {
	Key   string
	Value string
}

// size and server ip mapping
type AllocationBlock struct {
	//Ex, http://127.0.0.1:123
	Server_url string

	// Value Size
	Size string
}
