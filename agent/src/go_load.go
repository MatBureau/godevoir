package main

import (
	"client/moninfluxdb"
	"log"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/load"
)

func goLoad(client *influxdb3.Client) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		mstats, err := load.Avg()
		if err != nil {
			log.Printf("Erreur dans le load")
		}
		Datas.Load = mstats
		percs, err := cpu.Percent(1*time.Second, true) // une valeur par CPU logique
		if err != nil {
			return
		}
		Datas.CPULoad = &percs
		err = moninfluxdb.WriteLoad(client, ServerURL+"/load")
		if err != nil {
			LogMessage("erreur : Base influx non joignable")
		}
		LogMessage("goroutine: goLoad")
	}
}
