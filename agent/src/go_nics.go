package main

import (
	"client/moninfluxdb"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
)

func goNics(client *influxdb3.Client) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		// fenêtre de mesure : 1 seconde (compromis lisibilité/latence)
		rates, err := NICRates(1 * time.Second)
		if err != nil {
			continue
		}
		Datas.Nics = &rates
		_ = moninfluxdb.WriteNics(client, ServerURL+"/nics")

		LogMessage("goroutine: goNics")
	}
}
