package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func webcpubyid(w http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Fprintf(w, "Erreur dans le cpu")
		return
	}

	mcpus := *Datas.CPU
	j, _ := json.Marshal(mcpus[id])
	// Active CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	log.Println("/cpu/ID")

}

func webcpu(w http.ResponseWriter, req *http.Request) {
	mcpus := *Datas.CPU

	j, _ := json.Marshal(mcpus)
	// Active CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	log.Println("/cpu")
}

func webcpuload(w http.ResponseWriter, req *http.Request) {

	j, _ := json.Marshal(Datas.CPULoad)
	// Active CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	log.Println("/cpu/load")
}
