package main

import (
	"log"
	"./drone"
	"image"
	"./decoder"
	"./ui"
)

var videoChannel = make(chan *image.Image)
var commandChannel = make(chan interface{})
var flightData = make(chan string)

func main() {
	err := decoder.Init()
	if err != nil {
		log.Fatal("Unable to create decoder", err)
	}

	defer decoder.Free()

	go drone.DroneControl(videoChannel, commandChannel, flightData)

	ui.Start(videoChannel, commandChannel, flightData)
}
