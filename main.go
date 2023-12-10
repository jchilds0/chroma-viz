package main

import (
	"chroma-viz/gui"
)

func main() {
    gui.InitConnections()
    gui.AddConnection("Engine", "127.0.0.1", 6800)
    gui.AddConnection("Preview", "127.0.0.1", 6100)

    gui.SetupMainGui()
}

