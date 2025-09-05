package main

import (
	"fmt"
	"net/http"
)

const ServerURL string = "http://localhost:8080"
const DEBUG bool = false

// Adresses des agents distants
var AgentHosts = []string{
	"162.19.228.206:8080",
	"10.4.249.146:8080",
	"192.168.89.132:8080",
}

//var client *influxdb3.Client

func main() {
	http.HandleFunc("GET /html/nics", htmlnics)
	http.HandleFunc("GET /html/disks", htmldisks)
	http.HandleFunc("GET /html/load", htmlload)
	http.HandleFunc("GET /html/procs", htmlprocs)
	http.HandleFunc("GET /html/cpus", htmlcpus)
	http.HandleFunc("GET /html/memory", htmlmem)
	http.HandleFunc("GET /html/agents", htmlagents)
	http.HandleFunc("GET /api/agents", apiagentstatus)
	http.HandleFunc("GET /api/procs", apiprocs)
	http.HandleFunc("POST /api/procs/kill/{pid}", apiprockill)
	fmt.Println("Serveur :9090")
	http.ListenAndServe(":9090", nil)
}
