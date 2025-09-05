package main

import (
	"client/moninfluxdb"
	"log"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
)

func goProcs(client *influxdb3.Client) {
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
		_ = moninfluxdb.WriteProcsCount(client, ServerURL+"/procs")

		LogMessage("goroutine: goProcs")
	}
}
