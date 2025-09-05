package main

import "log"

func LogMessage(message string) {
	if DEBUG {
		log.Println(message)
	}
}
