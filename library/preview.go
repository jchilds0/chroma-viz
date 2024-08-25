package library

import (
	"log"
	"os/exec"
	"strconv"

	"github.com/gotk3/gotk3/gtk"
)

func SetupPreviewWindow(conf Config, takeOn, cont, takeOff func()) (box *gtk.Box, err error) {
	box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}

	padding := uint(10)
	actions, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}

	box.PackStart(actions, false, false, padding)

	takeOnButton, err := gtk.ButtonNewWithLabel("Play")
	if err != nil {
		return
	}

	takeOnButton.Connect("clicked", takeOn)
	actions.PackStart(takeOnButton, false, false, padding)

	contButton, err := gtk.ButtonNewWithLabel("Continue")
	if err != nil {
		return
	}

	contButton.Connect("clicked", cont)
	actions.PackStart(contButton, false, false, padding)

	takeOffButton, err := gtk.ButtonNewWithLabel("Play Off")
	if err != nil {
		return
	}

	takeOffButton.Connect("clicked", takeOff)
	actions.PackStart(takeOffButton, false, false, padding)

	restart, err := gtk.ButtonNewWithLabel("Restart Preview")
	if err != nil {
		return
	}

	actions.PackEnd(restart, false, false, padding)

	window, err := gtk.FrameNew("")
	if err != nil {
		return
	}

	box.PackStart(window, true, true, 0)
	window.Connect("draw", func(window *gtk.Frame) {
		width := window.GetAllocatedWidth()
		height := width * 9 / 16
		window.SetSizeRequest(-1, height)
	})

	window.SetVisible(true)
	soc, err := startPreview(conf)
	if err != nil {
		return
	}

	window.Add(soc)

	restart.Connect("clicked", func() {
		soc.Destroy()

		soc, err = startPreview(conf)
		if err != nil {
			log.Print(err)
			return
		}

		window.Add(soc)
	})

	return
}

func startPreview(conf Config) (soc *gtk.Socket, err error) {
	soc, err = gtk.SocketNew()
	if err != nil {
		return
	}

	soc.SetVisible(true)
	soc.Connect("realize", func(soc *gtk.Socket) {
		xid := soc.GetId()

		chromaEnginePath := conf.PreviewDirectory + conf.PreviewName
		prev := exec.Command(
			chromaEnginePath,
			"-w", strconv.Itoa(int(xid)),
			"-c", conf.PreviewConfig,
		)

		log.Print(prev.String())

		if err := prev.Start(); err != nil {
			log.Print(err)
		}
	})

	soc.Connect("plug-added", func() { log.Printf("Preview window connected") })
	return
}
