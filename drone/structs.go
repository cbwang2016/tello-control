package drone

type TakeOffCommand struct {

}

type ThrowTakeOffCommand struct {

}

type LandCommand struct {
  
}

type StopLandingCommand struct {
  
}

type SetFastModeCommand struct {
	
}

type SetSlowModeCommand struct {
	
}

type RotateClockwiseCommand struct {
	Value int
}

type RotateCounterClockwiseCommand struct {
	Value int
}

type UpCommand struct {
	Value int
}

type DownCommand struct {
	Value int
}

type LeftCommand struct {
	Value int
}

type RightCommand struct {
	Value int
}

type ForwardCommand struct {
	Value int
}

type BackwardCommand struct {
	Value int
}

type FlipForwardCommand struct {

}

type FlipBackwardCommand struct {

}

type FlipLeftCommand struct {

}

type FlipRightCommand struct {

}
