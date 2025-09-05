package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func webnics(w http.ResponseWriter, req *http.Request) {
	// Active CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(Datas.Nics)
	log.Println("/nics")
}
