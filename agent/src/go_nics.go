package main

import (
	"time"
)

func goNics() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		// fenêtre de mesure : 1 seconde (compromis lisibilité/latence)
		rates, err := NICRates(1 * time.Second)
		if err != nil {
			continue
		}
		Datas.Nics = &rates
		LogMessage("goroutine: goNics")
	}
}
