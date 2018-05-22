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
	var fd tello.FlightData
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
			drone.SetVideoEncoderRate(0)
			gobot.Every(500*time.Millisecond, func() {
				drone.StartVideo()
			})
		})

		drone.On(tello.LightStrengthEvent, func(data interface{}) {
			fd.LightStrength = data.(int16)
		})

		drone.On(tello.WifiDataEvent, func(data interface{}) {
			tmp := data.(*tello.WifiData)
			fd.WifiStrength = tmp.Strength
			fd.WifiDisturb = tmp.Disturb
		})

		drone.On(tello.FlightDataEvent, func(data interface{}) {
			fd2 := data.(*tello.FlightData)
			fd.BatteryPercentage = fd2.BatteryPercentage
			fd.Height = fd2.Height
			fd.NorthSpeed = fd2.NorthSpeed
			fd.EastSpeed = fd2.EastSpeed
			fd.DroneHover = fd2.DroneHover
			fd.EmSky = fd2.EmSky
			fd.EmGround = fd2.EmGround
			fd.EmOpen = fd2.EmOpen
			localFlightData = ""
			if (fd2.BatteryLow) {
				localFlightData += "Warning: BatteryLow! "
			}
			if (fd2.BatteryLower) {
				localFlightData += "Warning: BatteryLower! "
			}
			if (fd2.BatteryState) {
				localFlightData += "Warning: BatteryState! "
			}
			if (fd2.DownVisualState) {
				localFlightData += "Warning: DownVisualState! "
			}
			if (fd2.GravityState) {
				localFlightData += "Warning: GravityState! "
			}
			if (fd2.ImuState) {
				localFlightData += "Warning: ImuState! "
			}
			if (fd2.PowerState) {
				localFlightData += "Warning: PowerState! "
			}
			if (fd2.PressureState) {
				localFlightData += "Warning: PressureState! "
			}
			if (fd2.WindState) {
				localFlightData += "Warning: WindState! "
			}
			if (fd2.WifiDisturb != 0) {
				localFlightData += "Warning: WifiDisturb! "
			}
			localFlightData += fmt.Sprintf("Batt: %d%%, WifiStrength: %d, WifiDisturb: %d, Height: %.1fm, NorthSpeed: %.1fm/s, EastSpeed: %.1fm/s, LightStrength: %d",
				fd.BatteryPercentage, fd.WifiStrength, fd.WifiDisturb,
				float32(fd.Height)/10, float32(fd.NorthSpeed)/10, float32(fd.EastSpeed)/10,
				fd.LightStrength)
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
				case ThrowTakeOffCommand:
					log.Printf("Going to (Throw & Go)take off, %q", cmd)
					drone.ThrowTakeOff()
				case LandCommand:
					log.Printf("Going to land, %q", cmd)
					drone.Land()
				case StopLandingCommand:
					log.Printf("Stop landing, %q", cmd)
					drone.StopLanding()
				case SetFastModeCommand:
					log.Printf("Setting Fast Mode, %q", cmd)
					drone.SetFastMode()
				case SetSlowModeCommand:
					log.Printf("Setting Slow Mode, %q", cmd)
					drone.SetSlowMode()
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

