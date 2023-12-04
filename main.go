package main

import (
	"chroma-viz/gui"
)

func main() {
    conn := gui.NewConnection("127.0.0.1")
    gui.LaunchGui(conn)

    //conn.CloseConn()

    // if prev != nil {
    //     prev.Wait()
    // }
}

