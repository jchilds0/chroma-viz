package gui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

const (
    engineWidth = 10
    engineHeight = 10
)

type EngineWidget struct {
    *gtk.Button
    box          *gtk.Box
    area         *gtk.DrawingArea
    conn         *Connection
    connStatus   bool
}

func NewEngineWidget(name string, conn *Connection) *EngineWidget {
    eng := &EngineWidget{conn: conn, connStatus: false}
    eng.Button, _ = gtk.ButtonNew()
    eng.box, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    eng.Button.Add(eng.box)

    label, _ := gtk.LabelNew(name + " ")
    eng.box.PackStart(label, true, true, 0)

    eng.area, _ = gtk.DrawingAreaNew()
    eng.box.PackStart(eng.area, true, true, 0)

    eng.Connect(
        "clicked", 
        func() { 
            if eng.connStatus {
                eng.connStatus = eng.conn.IsConnected()
            } else {
                eng.connStatus = eng.conn.Connect()
            } 

            eng.area.QueueDraw()
        })

    eng.area.Connect(
        "draw", 
        func(da *gtk.DrawingArea, cr *cairo.Context) {
            height := da.GetAllocatedHeight()
            da.SetSizeRequest(height, height)

            if eng.connStatus {
                cr.SetSourceRGB(0, 255, 0)
            } else {
                cr.SetSourceRGB(255, 0, 0)
            }
            
            cr.Rectangle(0, 0, float64(height), float64(height))
            cr.Fill()
        })

    eng.Button.QueueDraw()

    return eng
}

