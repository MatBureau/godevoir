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
	// Routes pour les pages HTML
	http.HandleFunc("GET /html/index", htmlindex)
	http.HandleFunc("GET /html/nics", htmlnics)
	http.HandleFunc("GET /html/disks", htmldisks)
	http.HandleFunc("GET /html/load", htmlload)
	http.HandleFunc("GET /html/procs", htmlprocs)
	http.HandleFunc("GET /html/cpus", htmlcpus)
	http.HandleFunc("GET /html/memory", htmlmem)
	http.HandleFunc("GET /html/agents", htmlagents)
	
	// Routes API
	http.HandleFunc("GET /api/agents", apiagentstatus)
	http.HandleFunc("GET /api/agent/data", apiagentdata)
	
	fmt.Println("Serveur :9090")
	http.ListenAndServe(":9090", nil)
}
