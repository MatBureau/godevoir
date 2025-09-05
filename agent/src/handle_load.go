package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func webload(w http.ResponseWriter, req *http.Request) {
	j, _ := json.Marshal(Datas.Load)
	// Active CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	log.Println("/load")
}
