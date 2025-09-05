package main

import (
	"client/moninfluxdb"
	"log"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/shirou/gopsutil/v4/cpu"
)

func goCPU(client *influxdb3.Client) {
	ticker := time.NewTicker(800 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		mcpu, err := cpu.Info()
		Datas.CPU = &mcpu
		_ = moninfluxdb.WriteCPUPercent(client, ServerURL+"/cpu/load")
		if err != nil {
			log.Println("Erreur dans le cpu")
			return
		}

		time.Sleep(time.Millisecond * 800)
		LogMessage("goroutine: goCPU")
	}
}
