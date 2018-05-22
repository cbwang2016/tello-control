package ui

import (
	"github.com/hajimehoshi/ebiten"
	"log"
	"image"
	"github.com/hajimehoshi/ebiten/inpututil"
	"../drone"
	"math"
	
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
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

func Start(videoChannel chan *image.Image, commandChannel chan interface{}, flightData chan string) {
	var lastImage *ebiten.Image = nil
	//var tookOff = false
	var isFastMode = false
	var localFlightData string
	tt, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	uiFont := truetype.NewFace(tt, &truetype.Options{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	update := func (screen *ebiten.Image) error {
		for _, id := range inpututil.JustConnectedGamepadIDs() {
			log.Printf("gamepad connected: id: %d", id)
		}

		if(inpututil.IsKeyJustPressed(ebiten.KeyF)) {
			if (isFastMode) {
				commandChannel <- drone.SetSlowModeCommand{}
			} else {
				commandChannel <- drone.SetFastModeCommand{}
			}
			isFastMode = !isFastMode
		}

		if(inpututil.IsKeyJustPressed(ebiten.KeyA)) {
			commandChannel <- drone.RotateCounterClockwiseCommand{100}
		}
		if(inpututil.IsKeyJustPressed(ebiten.KeyD)) {
			commandChannel <- drone.RotateClockwiseCommand{100}
		}
		if(inpututil.IsKeyJustReleased(ebiten.KeyA) || inpututil.IsKeyJustReleased(ebiten.KeyD)) {
			commandChannel <- drone.RotateClockwiseCommand{0}
		}
		
		if(inpututil.IsKeyJustPressed(ebiten.KeyW)) {
			commandChannel <- drone.UpCommand{100}
		}
		if(inpututil.IsKeyJustPressed(ebiten.KeyS)) {
			commandChannel <- drone.DownCommand{100}
		}
		if(inpututil.IsKeyJustReleased(ebiten.KeyW) || inpututil.IsKeyJustReleased(ebiten.KeyS)) {
			commandChannel <- drone.UpCommand{0}
		}

		if(inpututil.IsKeyJustPressed(ebiten.KeyUp)) {
			if inpututil.KeyPressDuration(ebiten.KeyControl) > 0 {
				commandChannel <- drone.FlipForwardCommand{}
			}
			commandChannel <- drone.ForwardCommand{100}
		}
		if(inpututil.IsKeyJustPressed(ebiten.KeyDown)) {
			if inpututil.KeyPressDuration(ebiten.KeyControl) > 0 {
				commandChannel <- drone.FlipBackwardCommand{}
			}
			commandChannel <- drone.BackwardCommand{100}
		}
		if(inpututil.IsKeyJustReleased(ebiten.KeyDown) || inpututil.IsKeyJustReleased(ebiten.KeyUp)) {
			commandChannel <- drone.ForwardCommand{0}
		}
		
		if(inpututil.IsKeyJustPressed(ebiten.KeyLeft)) {
			if inpututil.KeyPressDuration(ebiten.KeyControl) > 0 {
				commandChannel <- drone.FlipLeftCommand{}
			}
			commandChannel <- drone.LeftCommand{100}
		}
		if(inpututil.IsKeyJustPressed(ebiten.KeyRight)) {
			if inpututil.KeyPressDuration(ebiten.KeyControl) > 0 {
				commandChannel <- drone.FlipRightCommand{}
			}
			commandChannel <- drone.RightCommand{100}
		}
		if(inpututil.IsKeyJustReleased(ebiten.KeyLeft) || inpututil.IsKeyJustReleased(ebiten.KeyRight)) {
			commandChannel <- drone.LeftCommand{0}
		}
		
		if(inpututil.IsKeyJustReleased(ebiten.KeyT)) {
			if inpututil.KeyPressDuration(ebiten.KeyControl) > 0 {
				commandChannel <- drone.ThrowTakeOffCommand{}
			} else {
				commandChannel <- drone.TakeOffCommand{}
			}
		}
		if(inpututil.IsKeyJustReleased(ebiten.KeyL)) {
			commandChannel <- drone.LandCommand{}
		}
		if(inpututil.IsKeyJustReleased(ebiten.KeyC)) {
			commandChannel <- drone.StopLandingCommand{}
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
			case localFlightData = <-flightData:

			default:
		}

		if lastImage != nil {
			screen.DrawImage(lastImage, nil)
			if isFastMode {
				text.Draw(screen, localFlightData + ", Mode: Fast", uiFont, 20, 20, color.White)
			} else {
				text.Draw(screen, localFlightData + ", Mode: Slow", uiFont, 20, 20, color.White)
			}
		}

		return nil
	}

	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Drone Control"); err != nil {
		log.Fatal(err)
	}
}
