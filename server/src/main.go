package main

import (
	"fmt"
	"net/http"
)

const ServerURL string = "http://localhost:8080"
const DEBUG bool = false

//var client *influxdb3.Client

func main() {
	http.HandleFunc("GET /html/nics", htmlnics)
	http.HandleFunc("GET /html/disks", htmldisks)
	http.HandleFunc("GET /html/load", htmlload)
	http.HandleFunc("GET /html/procs", htmlprocs)
	http.HandleFunc("GET /html/cpus", htmlcpus)
	http.HandleFunc("GET /html/memory", htmlmem)
	fmt.Println("Serveur :9090")
	http.ListenAndServe(":9090", nil)
}
