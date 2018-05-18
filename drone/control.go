package drone

import (
	"log"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"image"
	"fmt"
	"time"
	"../decoder"
	"os"
)

func DroneControl(videoChannel chan *image.Image, commandChannel chan interface{}, flightData chan string) {
	var localFlightData string
	var fd *tello.FlightData
	os.MkdirAll("recordings", os.ModePerm)
	t := time.Now()

	f, err := os.Create("recordings/" + t.Format("2006-01-02T15-04-05.nal"))
	if err != nil {
		panic("Unable to create recording file")
	}

	defer f.Close()

	drone := tello.NewDriver("8890")

	imageHandler := func(im *image.Image) {
		videoChannel <- im
	}

	work := func() {
		drone.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Connected")
			drone.StartVideo()
			drone.SetVideoEncoderRate(5)
			gobot.Every(100*time.Millisecond, func() {
				drone.StartVideo()
			})
		})

		drone.On(tello.LightStrengthEvent, func(data interface{}) {
			fd.LightStrength = data.(int16)
		})

		drone.On(tello.FlightDataEvent, func(data interface{}) {
			fd = data.(*tello.FlightData)
			localFlightData = ""
			if (fd.BatteryLow) {
				localFlightData += "Warning: BatteryLow!"
			}
			if (fd.BatteryLower) {
				localFlightData += "Warning: BatteryLower!"
			}
			if (fd.BatteryState) {
				localFlightData += "Warning: BatteryState!"
			}
			if (fd.DownVisualState) {
				localFlightData += "Warning: DownVisualState!"
			}
			if (fd.GravityState) {
				localFlightData += "Warning: GravityState!"
			}
			if (fd.ImuState) {
				localFlightData += "Warning: ImuState!"
			}
			if (fd.PowerState) {
				localFlightData += "Warning: PowerState!"
			}
			if (fd.PressureState) {
				localFlightData += "Warning: PressureState!"
			}
			if (fd.WindState) {
				localFlightData += "Warning: WindState!"
			}
			localFlightData += fmt.Sprintf("Batt: %d%%, WifiStrength: %d, Height: %.1fm, Speed: %.1fm/s, Hover: %t, Sky: %t, Ground: %t, Open: %t, LightStrength: %d",
				fd.BatteryPercentage, fd.WifiStrength,
				float32(fd.Height)/10, float32(fd.FlySpeed)/10,
				fd.DroneHover,
				fd.EmSky, fd.EmGround, fd.EmOpen, fd.LightStrength)
			flightData <- localFlightData
		})


		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			decoder.Decode(pkt, imageHandler)

			// dump NALs
			_, err := f.Write(pkt)
			if err != nil {
				panic("Unable to write recording")
			}

			f.Sync()
		})

		for {
			select {
			case cmd := <-commandChannel:
				switch cmd := cmd.(type) {
				case TakeOffCommand:
					log.Printf("Going to take off, %q", cmd)
					drone.TakeOff()
				case LandCommand:
					log.Printf("Going to land, %q", cmd)
					drone.Land()
				case RotateCounterClockwiseCommand:
					log.Printf("Rotating counter-clockwise %d", cmd.Value)
					drone.CounterClockwise(cmd.Value)
				case RotateClockwiseCommand:
					log.Printf("Rotating clockwise %d", cmd.Value)
					drone.Clockwise(cmd.Value)
				case UpCommand:
					log.Printf("Going up %d", cmd.Value)
					drone.Up(cmd.Value)
				case DownCommand:
					log.Printf("Going down %d", cmd.Value)
					drone.Down(cmd.Value)
				case LeftCommand:
					log.Printf("Going left %d", cmd.Value)
					drone.Left(cmd.Value)
				case RightCommand:
					log.Printf("Going right %d", cmd.Value)
					drone.Right(cmd.Value)
				case ForwardCommand:
					log.Printf("Going forward %d", cmd.Value)
					drone.Forward(cmd.Value)
				case BackwardCommand:
					log.Printf("Going backward %d", cmd.Value)
					drone.Backward(cmd.Value)
				case FlipForwardCommand:
					log.Printf("Front Flip")
					drone.FrontFlip()
				case FlipBackwardCommand:
					log.Printf("Back Flip")
					drone.BackFlip()
				case FlipLeftCommand:
					log.Printf("Left Flip")
					drone.LeftFlip()
				case FlipRightCommand:
					log.Printf("Right Flip")
					drone.RightFlip()
				}
			}
		}
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}

