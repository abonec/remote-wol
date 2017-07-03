package main

import (
	"github.com/mrmorphic/hwio"
)

func turnRelay() {
	relayPin, err := hwio.GetPinWithMode("gpio17", hwio.OUTPUT)
	defer hwio.UnassignPin(relayPin)
	defer hwio.ClosePin(relayPin)
	failError(err)

	hwio.DigitalWrite(relayPin, hwio.HIGH)
	hwio.Delay(300)
	hwio.DigitalWrite(relayPin, hwio.LOW)
	defer hwio.UnassignPin(relayPin)
}
