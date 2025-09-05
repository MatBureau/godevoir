package main

import (
	"encoding/json"
	"net/http"
)

func webdiskusage(w http.ResponseWriter, r *http.Request) {
	j, _ := json.Marshal(Datas.DiskUsage)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
