package main

import (
	"github.com/abonec/go-ping"
	"time"
)

func startChecker() {
	isOnline := pingMachine()
	for {
		isCurrentOnline := pingMachine()
		if isCurrentOnline != isOnline {
			isOnline = isCurrentOnline
			sendStatus(isOnline)
		}
		time.Sleep(5*time.Second)
	}
}

func pingMachine() bool {
	pinger, err := ping.NewPinger(defaultIp)
	printError(err)
	pinger.Count = 1
	pinger.Timeout = time.Second
	pinger.Run()
	if pinger.Statistics().PacketsRecv > 0 {
		return true
	} else {
		return false
	}
}

func sendStatus(isOnline bool) {
	if isOnline {
		sendGroupMessage("machine is online now")
	} else {
		sendGroupMessage("machine is offline now")
	}

}

