package main

import (
	"chroma-viz/gui"
)

func main() {
    conn := make(map[string]*gui.Connection)
    conn["Engine"] = gui.NewConnection("127.0.0.1", 6800)
    conn["Preview"] = gui.NewConnection("127.0.0.1", 6100)
    gui.LaunchGui(conn)
}

