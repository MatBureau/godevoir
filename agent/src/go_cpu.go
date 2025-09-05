package main

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

func goCPU() {
	ticker := time.NewTicker(800 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		mcpu, err := cpu.Info()
		Datas.CPU = &mcpu
		if err != nil {
			log.Println("Erreur dans le cpu")
			return
		}

		time.Sleep(time.Millisecond * 800)
		LogMessage("goroutine: goCPU")
	}
}
