package main

import (
	"client/moninfluxdb"
	"fmt"
	"net/http"
)

const DBHost string = "http://localhost:8181"
const DBName string = "load"
const DBToken string = "apiv3_9rdZKWDNT9DcuIli9W6R9ePr9KpeRXTkWMktG_1B5r504DRNMsYviUUz5107Hq1K9C4xmIvJDo7YayM1nd_Wsg"
const ServerURL string = "http://162.19.228.206:8080"
const DEBUG bool = false

func main() {
	// Connexion InfluxDB
	client, err := moninfluxdb.Open(DBHost, DBName, DBToken)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// MaJ de AllDatas
	go goLoad(client)
	go goCPU(client)
	go goDisk()
	go goProcs(client)
	go goNics(client)
	go goMem(client)
	go goDiskUsage(client)

	http.HandleFunc("GET /cpu", webcpu)
	http.HandleFunc("GET /cpu/{id}", webcpubyid)
	http.HandleFunc("GET /cpu/load", webcpuload)
	http.HandleFunc("GET /load", webload)
	http.HandleFunc("GET /procs", webprocs)
	http.HandleFunc("GET /procs/{id}", webprocsbypid)
	http.HandleFunc("GET /disks", webdisk)
	http.HandleFunc("GET /nics", webnics)
	http.HandleFunc("GET /mem", webmem)
	http.HandleFunc("GET /procs/kill/{pid}", webprocskill)
	http.HandleFunc("GET /disks/usage", webdiskusage)
	fmt.Println("Serveur :8080")
	http.ListenAndServe(":8080", nil)
}
