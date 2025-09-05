package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"syscall"
)

type killResp struct {
	PID    int    `json:"pid"`
	Action string `json:"action"`
	Ok     bool   `json:"ok"`
	Error  string `json:"error,omitempty"`
}

func webprocskill(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	pidStr := req.PathValue("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid <= 0 {
		http.Error(w, "pid invalide", http.StatusBadRequest)
		return
	}

	p, err := os.FindProcess(pid)
	if err != nil {
		http.Error(w, "process introuvable", http.StatusNotFound)
		return
	}

	// D’abord TERM, puis KILL si nécessaire (Unix) ; Windows: Kill().
	fmt.Println("OS détecté :", (runtime.GOOS))
	kill := func() error {
		if runtime.GOOS == "windows" {
			return p.Kill()
		}
		if err := p.Signal(syscall.SIGTERM); err != nil {
			return err
		}
		if err := p.Signal(syscall.SIGKILL); err != nil {
			return err
		}
		return nil
	}

	resp := killResp{PID: pid, Action: "kill"}
	if err := kill(); err != nil {
		resp.Ok = false
		resp.Error = err.Error()
		w.WriteHeader(http.StatusBadGateway)
	} else {
		resp.Ok = true
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func webprocs(w http.ResponseWriter, req *http.Request) {

	j, _ := json.Marshal(Datas.Procs)
	// Active CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	log.Println("/procs")
}

func webprocsbypid(w http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	idint, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Fprintf(w, "Erreur dans le processus")
		return
	}
	id := int32(idint)
	out, err := DTOProcLoad(id)
	if err != nil {
		fmt.Fprintf(w, "Erreur dans le processus")
		return
	}
	j, _ := json.Marshal(out)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	log.Println("/procs/ID")
}
