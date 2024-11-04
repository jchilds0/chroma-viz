package library

// #cgo pkg-config: gtk+-3.0 glew freetype2
// #cgo CFLAGS: -I/home/josh/programming/chroma-engine/src
// #cgo LDFLAGS: -L/home/josh/programming/chroma-engine/build -lchroma -lm -lpng
// #include "chroma-typedefs.h"
// #include <stdlib.h>
import "C"

import (
	"errors"
	"log"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
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
	glArea, err := startPreview(conf)
	if err != nil {
		return
	}

	window.Add(glArea)

	restart.Connect("clicked", func() {
		log.Println("Restart Preview not implemented")
		return

		glArea, err = startPreview(conf)
		if err != nil {
			log.Print(err)
			return
		}

		window.Add(glArea)
	})

	return
}

func startPreview(conf Config) (prev *gtk.GLArea, err error) {
	confStr := C.CString(conf.PreviewConfig)
	defer C.free(unsafe.Pointer(confStr))

	status := C.chroma_init_renderer(confStr)
	if status < 0 {
		err = errors.New("Error initializing preview renderer")
		return
	}

	c := C.chroma_new_renderer()
	if c == nil {
		err = errors.New("cgo returned unexpected nil pointer")
		return
	}

	obj := glib.Take(unsafe.Pointer(c))
	prev = &gtk.GLArea{
		Widget: gtk.Widget{
			InitiallyUnowned: glib.InitiallyUnowned{
				Object: obj,
			},
		},
	}

	return
}
