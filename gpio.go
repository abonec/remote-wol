package main

import (
	"github.com/mrmorphic/hwio"
	"sync"
	"log"
)

const (
	relayPinName = "gpio17"
)

var (
	mutex = &sync.Mutex{}
)

func turnRelay() {
	mutex.Lock()
	defer mutex.Unlock()
	defer hwio.CloseAll()
	relayPin, err := hwio.GetPinWithMode(relayPinName, hwio.OUTPUT)
	if err != nil {
		log.Println("Can't use relay")
		return
	}
	defer hwio.UnassignPin(relayPin)
	defer hwio.ClosePin(relayPin)
	defer hwio.UnassignPin(relayPin)

	hwio.DigitalWrite(relayPin, hwio.HIGH)
	hwio.Delay(300)
	hwio.DigitalWrite(relayPin, hwio.LOW)
}
