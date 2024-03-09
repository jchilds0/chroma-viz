package props

import (
	"chroma-viz/attribute"
	"log"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

type ClockEditor struct {
    box             *gtk.Box
    c               chan int
    timeFormat      string 
    edit            map[string]attribute.Editor
}

func NewClockEditor(width, height int, animate, cont func()) (clockEdit *ClockEditor, err error) {
    clockEdit = &ClockEditor{timeFormat: "04:05"}

    clockEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    clockEdit.box.SetVisible(true)
    clockEdit.c = make(chan int, 1)

    actions, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    clockEdit.box.PackStart(actions, false, false, padding)
    actions.SetVisible(true)

    startButton, err := gtk.ButtonNewWithLabel("Start")
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    actions.PackStart(startButton, false, false, padding)
    startButton.SetVisible(true)
    startButton.Connect("clicked", func() {
        clockEdit.c <- START
    })

    pauseButton, err := gtk.ButtonNewWithLabel("Pause")
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    actions.PackStart(pauseButton, false, false, padding)
    pauseButton.SetVisible(true)
    pauseButton.Connect("clicked", func() {
        clockEdit.c <- PAUSE 
    })

    stopButton, err := gtk.ButtonNewWithLabel("Stop")
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    actions.PackStart(stopButton, false, false, padding)
    stopButton.SetVisible(true)
    stopButton.Connect("clicked", func() {
        clockEdit.c <- STOP
    })

    clockEdit.edit = make(map[string]attribute.Editor, 5)
    clockEdit.edit["string"], err = attribute.NewStringEditor("Time", animate)
    if err != nil {
        return 
    }

    clockEdit.edit["x"], err = attribute.NewIntEditor("x", 0, float64(width), animate)
    if err != nil {
        return
    }

    clockEdit.edit["y"], err = attribute.NewIntEditor("y", 0, float64(height), animate)
    if err != nil {
        return
    }

    clockEdit.edit["color"], err = attribute.NewColorEditor("Color", animate)
    if err != nil {
        return
    }

    clockEdit.box.PackStart(clockEdit.edit["x"].Box(), false, false, padding)
    clockEdit.box.PackStart(clockEdit.edit["y"].Box(), false, false, padding)
    clockEdit.box.PackStart(clockEdit.edit["string"].Box(), false, false, padding)
    clockEdit.box.PackStart(clockEdit.edit["color"].Box(), false, false, padding)

    go clockEdit.RunClock(cont)
    return 
}

func (clock *ClockEditor) RunClock(cont func()) {
    state := PAUSE
    tick := time.NewTicker(time.Second)

    timeEdit := clock.edit["string"]
    if timeEdit == nil {
        log.Fatalf("Missing time editor in clock")
    } 

    timeEditor, ok := timeEdit.(*attribute.StringEditor)
    if !ok {
        log.Fatalf("time edit is not a StringEditor")
    }

    for {
        select {
        case state = <-clock.c:
        case <-tick.C:
        }    

        switch state {
        case START:
            // update time and animate
            currentText, err := timeEditor.Entry.GetText()
            if err != nil {
                log.Println(err)
                continue
            }

            currentTime, err := time.Parse(clock.timeFormat, currentText)
            if err != nil {
                log.Println(err)
                continue
            }

            currentTime = currentTime.Add(time.Second)
            timeEditor.Entry.SetText(currentTime.Format(clock.timeFormat))
            tick = time.NewTicker(time.Second)

            cont()
        case PAUSE:
            // block until we recieve an instruction
            state = <-clock.c
        case STOP:
            // reset the time and block
            timeEditor.Entry.SetText("00:00")
            state = <-clock.c
        default:
            log.Printf("Clock recieved unknown value through channel %d\n", state)
        }
    }
}

func (clock *ClockEditor) Box() *gtk.Box {
    return clock.box
}

func (clock *ClockEditor) Editors() map[string]attribute.Editor {
    return clock.edit
}

