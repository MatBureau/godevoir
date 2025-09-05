package main

import (
	"log"
	"time"
)

func goProcs() {
	ticker := time.NewTicker(2000 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		// Remonte tous les processus
		out, err := DTOProcAllLoad()
		if err != nil {
			log.Println("Erreur dans le procs")
			return
		}
		Datas.Procs = out
		LogMessage("goroutine: goProcs")
	}
}
