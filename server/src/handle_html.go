package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
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
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := map[string]any{
		"ServerURL": ServerURL,
	}
	_ = tpl.Execute(w, data)
	log.Println("/html/procs")
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
