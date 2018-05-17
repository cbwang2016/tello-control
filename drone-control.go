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

func main() {
	err := decoder.Init()
	if err != nil {
		log.Fatal("Unable to create decoder", err)
	}

	defer decoder.Free()

	go drone.DroneControl(videoChannel, commandChannel)

	ui.Start(videoChannel, commandChannel)
}
