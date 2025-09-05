package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func htmlnics(w http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles("www/nics.html")
	if err != nil {
		fmt.Fprintf(w, "parse nics.html: %v", err)
	}

	//Ajour du type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := map[string]any{
		"ServerURL": ServerURL,
	}
	_ = tpl.Execute(w, data)
	log.Println("/html/nics")
}

func htmlprocs(w http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles("www/procs.html")
	if err != nil {
		fmt.Fprintf(w, "parse procs.html: %v", err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := map[string]any{
		"ServerURL":  ServerURL,          // si tu en as encore besoin ailleurs
		"ListURL":    "/api/procs",       // <- nouveau: proxy GET
		"KillPrefix": "/api/procs/kill/", // <- déjà utilisé pour le kill
	}
	_ = tpl.Execute(w, data)
	log.Println("/html/procs")
}

// GET /api/procs  ->  GET {ServerURL}/procs
func apiprocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	target := strings.TrimRight(ServerURL, "/") + "/procs"
	req, _ := http.NewRequest(http.MethodGet, target, nil)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func apiprockill(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pid := r.PathValue("pid")
	if _, err := strconv.Atoi(pid); err != nil {
		http.Error(w, "pid invalide", http.StatusBadRequest)
		return
	}
	target := strings.TrimRight(ServerURL, "/") + "/procs/kill/" + pid

	req, _ := http.NewRequest(http.MethodPost, target, nil)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func htmldisks(w http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles("www/disks.html")
	if err != nil {
		fmt.Fprintf(w, "parse disks.html: %v", err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := map[string]any{
		"ServerURL": ServerURL,
	}
	_ = tpl.Execute(w, data)
	log.Println("/html/disks")
}

func htmlload(w http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles("www/load.html")
	if err != nil {
		fmt.Fprintf(w, "parse load.html: %v", err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := map[string]any{
		"ServerURL": ServerURL,
	}
	_ = tpl.Execute(w, data)
	log.Println("/html/load")
}

func htmlcpus(w http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles("www/cpu.html")
	if err != nil {
		fmt.Fprintf(w, "parse cpu.html: %v", err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := map[string]any{
		"ServerURL": ServerURL,
	}
	_ = tpl.Execute(w, data)
	log.Println("/html/cpu")
}

func htmlmem(w http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles("www/mem.html")
	if err != nil {
		fmt.Fprintf(w, "parse cpu.html: %v", err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := map[string]any{
		"ServerURL": ServerURL,
	}
	_ = tpl.Execute(w, data)
	log.Println("/html/mem")
}

func htmlagents(w http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles("www/agents.html")
	if err != nil {
		fmt.Fprintf(w, "parse agents.html: %v", err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := map[string]any{
		"ServerURL": "http://localhost:9090",
	}
	_ = tpl.Execute(w, data)
	log.Println("/html/agents")
}

// Structure pour représenter le statut d'un agent
type AgentStatus struct {
	Host     string    `json:"host"`
	Status   string    `json:"status"`
	LastSeen time.Time `json:"last_seen"`
	CPU      []any     `json:"cpu,omitempty"`
	Memory   any       `json:"memory,omitempty"`
	Load     any       `json:"load,omitempty"`
	Procs    []any     `json:"procs,omitempty"`
	Error    string    `json:"error,omitempty"`
}

// Fonction pour vérifier le statut de tous les agents
func apiagentstatus(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")

	var agents []AgentStatus
	var wg sync.WaitGroup
	agentChan := make(chan AgentStatus, len(AgentHosts))

	// Vérifier chaque agent en parallèle
	for _, host := range AgentHosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			agent := checkAgentStatus(host)
			agentChan <- agent
		}(host)
	}

	// Fermer le canal une fois que tous les goroutines sont terminés
	go func() {
		wg.Wait()
		close(agentChan)
	}()

	// Collecter tous les résultats
	for agent := range agentChan {
		agents = append(agents, agent)
	}

	json.NewEncoder(w).Encode(agents)
}

// Fonction pour vérifier le statut d'un agent individuel
func checkAgentStatus(host string) AgentStatus {
	client := &http.Client{Timeout: 5 * time.Second}
	agent := AgentStatus{
		Host:     host,
		Status:   "offline",
		LastSeen: time.Now(),
	}

	// Tester la connexion avec /cpu
	resp, err := client.Get("http://" + host + "/cpu")
	if err != nil {
		agent.Error = err.Error()
		return agent
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		agent.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
		return agent
	}

	agent.Status = "online"

	// Récupérer les données CPU
	var cpu []any
	if err := json.NewDecoder(resp.Body).Decode(&cpu); err == nil {
		agent.CPU = cpu
	}

	// Récupérer les données mémoire
	if resp, err := client.Get("http://" + host + "/mem"); err == nil {
		defer resp.Body.Close()
		var mem any
		if err := json.NewDecoder(resp.Body).Decode(&mem); err == nil {
			agent.Memory = mem
		}
	}

	// Récupérer les données de charge
	if resp, err := client.Get("http://" + host + "/load"); err == nil {
		defer resp.Body.Close()
		var load any
		if err := json.NewDecoder(resp.Body).Decode(&load); err == nil {
			agent.Load = load
		}
	}

	// Récupérer les processus (limité aux 10 premiers pour éviter la surcharge)
	if resp, err := client.Get("http://" + host + "/procs"); err == nil {
		defer resp.Body.Close()
		var procs []any
		if err := json.NewDecoder(resp.Body).Decode(&procs); err == nil {
			// Limiter à 10 processus pour l'aperçu
			if len(procs) > 10 {
				procs = procs[:10]
			}
			agent.Procs = procs
		}
	}

	return agent
}
