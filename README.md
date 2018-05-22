# Tello Drone Control

Very early version with a lot of hacks. Based on https://github.com/socketbind/drone-control, heavily modified by cbwang2016.

**Needless to say I take no responsibility for any damage you caused to your drone. You should know exactly what you are doing before attempting to use any of this code.**

What works:
* Basic controls using a keyboard (up, down, rotate left, rotate right, forward, backward, left, right)
* Video stream
* Video save
* Flip controls work most of the time
* Flight data display(battery, height, speed, etc.)
![Preview](https://cbwang2016.github.io/images/Drone%20Control%202018_5_18%209_28_09.png)

## Usage

Fow Windows: download the exe file and three dlls. Remember to allow network access(Private & Public networks) in the Windows Firewall warning.

* Turn on the Tello
* Wait for it to initialise (flashing orange LED)
* Connect your computer to the Tello WiFi
* Run the exe

After a couple of seconds a video feed should appear - if it doesn't, then something is wrong so do not attempt to fly the Tello! You can try to rerun the exe if there's no video.

The recorded videos are in the "recordings/" folder. To convert it to a mp4 file, run:
```
ffmpeg -i "source.nal" -c:v copy -f mp4 "myOutputFile.mp4"
```

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

- T - Take off
- Ctrl+T - (Throw & Go)Take off
- L - Land
- C - Stop landing
- F - Fast/Slow mode switch

## Preqrequisites

* libavcodec - Used for decoding H.264 packets
* Gobot dev - dev branch of gobot