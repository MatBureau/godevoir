package main

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/disk"
)

func goDisk() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		parts, err := disk.Partitions(false)
		if err != nil {
			log.Println("Erreur dans les disques")
		}

		Datas.Parts = &parts
		LogMessage("goroutine: goDisk")
	}
}
