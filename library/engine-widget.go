package library

import (
	"log"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

const (
	engineWidth  = 10
	engineHeight = 10
)

type EngineWidget struct {
	Button *gtk.Button
	conn   *Connection
}

func NewEngineWidget(conn *Connection) *EngineWidget {
	var err error
	eng := &EngineWidget{conn: conn}

	eng.Button, err = gtk.ButtonNew()
	if err != nil {
		log.Fatalf("Error creating engine widget (%s)", err)
	}

	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		log.Fatalf("Error creating engine widget (%s)", err)
	}

	eng.Button.Add(box)

	label, err := gtk.LabelNew(conn.Name + " ")
	if err != nil {
		log.Fatalf("Error creating engine widget (%s)", err)
	}

	box.PackStart(label, true, true, 0)

	area, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatalf("Error creating engine widget (%s)", err)
	}

	box.PackStart(area, true, true, 0)

	eng.Button.Connect("clicked", func() {
		if !eng.conn.IsConnected() {
			eng.conn.Connect()
			go eng.conn.Watcher(func() { area.QueueDraw() })
		}
	})

	area.Connect("draw",
		func(da *gtk.DrawingArea, cr *cairo.Context) {
			height := da.GetAllocatedHeight()
			da.SetSizeRequest(height, height)

			if eng.conn.IsConnected() {
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
