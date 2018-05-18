# Tello Drone Control

Very early version with a lot of hacks. Based on https://github.com/socketbind/drone-control, heavily modified by cbwang2016.

**Needless to say I take no responsibility for any damage you caused to your drone. You should know exactly what you are doing before attempting to use any of this code.**

What works:
* Basic controls using a keyboard (up, down, rotate left, rotate right, forward, backward, left, right)
* Video stream
* Flip controls work most of the time
* Flight data display(battery, height, speed, etc.)
![Preview](https://cbwang2016.github.io/images/Drone%20Control%202018_5_18%209_10_40.png)

## Keyboard mappings

- W - Ascend
- S - Descend
- A - Rotate Counter-clockwise
- D - Rotate Clockwise

- Up - Forward
- Down - Backward
- Left - Sideways left
- Right - Sideways right

- Ctrl+Up - Flip forward
- Ctrl+Back - Flip backward
- Ctrl+Left - Flip left
- Ctrl+Right - Flip right

### Buttons
- T - Take off
- L - Land

## Preqrequisites

* libavcodec - Used for decoding H.264 packets
* Gobot dev - dev version of gobot

## Installation

Download the exe file and three dlls. Remember to allow Windows Firewall warning.