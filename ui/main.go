package ui

import (
	"github.com/hajimehoshi/ebiten"
	"log"
	"image"
	"github.com/hajimehoshi/ebiten/inpututil"
	"../drone"
	"math"
)

const (
	screenWidth  = 960
	screenHeight = 720

	takeOffButton = 1
	flipForwardButton = 12
	flipBackwardButton = 13
	flipLeftButton = 14
	flipRightButton = 15

	deadZoneHorizontal = 0.5
	deadZoneVertical = 0.5
)

func remapAxisInput(inputValue float64, deadZone float64, maxValue float64) int {
	if math.Abs(inputValue) > deadZone {
		if inputValue < 0 {
			return int(((inputValue + deadZone) / (1.0 - deadZone)) * maxValue)
		} else if inputValue > 0 {
			return int(((inputValue - deadZone) / (1.0 - deadZone)) * maxValue)
		}
	}

	return 0
}

func getV(fastMode bool) int {
	if(fastMode) {
		return 100
	} else {
		return 50
	}
}

func Start(videoChannel chan *image.Image, commandChannel chan interface{}) {
	var lastImage *ebiten.Image = nil
	//var tookOff = false
	var fastMode = false

	update := func (screen *ebiten.Image) error {
		for _, id := range inpututil.JustConnectedGamepadIDs() {
			log.Printf("gamepad connected: id: %d", id)
		}
		
		if(inpututil.IsKeyJustPressed(ebiten.KeyF)) {
			fastMode = !fastMode;
			log.Printf("Toggled fastMode. Now fastMode is %t", fastMode);
		}

		if(inpututil.IsKeyJustPressed(ebiten.KeyA)) {
			commandChannel <- drone.RotateCounterClockwiseCommand{getV(fastMode)}
		}
		if(inpututil.IsKeyJustPressed(ebiten.KeyD)) {
			commandChannel <- drone.RotateClockwiseCommand{getV(fastMode)}
		}
		if(inpututil.IsKeyJustReleased(ebiten.KeyA) || inpututil.IsKeyJustReleased(ebiten.KeyD)) {
			commandChannel <- drone.RotateClockwiseCommand{0}
		}
		
		if(inpututil.IsKeyJustPressed(ebiten.KeyW)) {
			commandChannel <- drone.UpCommand{getV(fastMode)}
		}
		if(inpututil.IsKeyJustPressed(ebiten.KeyS)) {
			commandChannel <- drone.DownCommand{getV(fastMode)}
		}
		if(inpututil.IsKeyJustReleased(ebiten.KeyW) || inpututil.IsKeyJustReleased(ebiten.KeyS)) {
			commandChannel <- drone.UpCommand{0}
		}

		if(inpututil.IsKeyJustPressed(ebiten.KeyUp)) {
			if inpututil.KeyPressDuration(ebiten.KeyControl) > 0 {
				commandChannel <- drone.FlipForwardCommand{}
			}
			commandChannel <- drone.ForwardCommand{getV(fastMode)}
		}
		if(inpututil.IsKeyJustPressed(ebiten.KeyDown)) {
			if inpututil.KeyPressDuration(ebiten.KeyControl) > 0 {
				commandChannel <- drone.FlipBackwardCommand{}
			}
			commandChannel <- drone.BackwardCommand{getV(fastMode)}
		}
		if(inpututil.IsKeyJustReleased(ebiten.KeyDown) || inpututil.IsKeyJustReleased(ebiten.KeyUp)) {
			commandChannel <- drone.ForwardCommand{0}
		}
		
		if(inpututil.IsKeyJustPressed(ebiten.KeyLeft)) {
			if inpututil.KeyPressDuration(ebiten.KeyControl) > 0 {
				commandChannel <- drone.FlipLeftCommand{}
			}
			commandChannel <- drone.LeftCommand{getV(fastMode)}
		}
		if(inpututil.IsKeyJustPressed(ebiten.KeyRight)) {
			if inpututil.KeyPressDuration(ebiten.KeyControl) > 0 {
				commandChannel <- drone.FlipRightCommand{}
			}
			commandChannel <- drone.RightCommand{getV(fastMode)}
		}
		if(inpututil.IsKeyJustReleased(ebiten.KeyLeft) || inpututil.IsKeyJustReleased(ebiten.KeyRight)) {
			commandChannel <- drone.LeftCommand{0}
		}
		
		if(inpututil.IsKeyJustReleased(ebiten.KeyT)) {
			commandChannel <- drone.TakeOffCommand{}
		}
		if(inpututil.IsKeyJustReleased(ebiten.KeyL)) {
			commandChannel <- drone.LandCommand{}
		}

		if ebiten.IsRunningSlowly() {
			return nil
		}

		select {
		case videoImage := <-videoChannel:
			var err error
			lastImage, err = ebiten.NewImageFromImage(*videoImage, ebiten.FilterDefault)
			if err != nil {
				panic("Unable to create image")
			}
		default:
		}

		if lastImage != nil {
			screen.DrawImage(lastImage, nil)
		}

		return nil
	}

	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Drone Control"); err != nil {
		log.Fatal(err)
	}
}
