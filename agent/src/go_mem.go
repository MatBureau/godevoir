package main

import (
	"client/moninfluxdb"
	"log"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/shirou/gopsutil/v4/mem"
)

func goMem(client *influxdb3.Client) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		v, err := mem.VirtualMemory()

		if err != nil {
			log.Println("Erreur dans la memoire")
		}

		Datas.Mem = v
		_ = moninfluxdb.WriteMemFromURL(client, ServerURL+"/mem", "hostA")
		LogMessage("goroutine: goMem")
	}
}
