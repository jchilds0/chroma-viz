package attribute

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

const (
    START = iota
    PAUSE
    STOP
)

type ClockAttribute struct {
    fileName    string
    chromaName  string 
    Value       string
}

func NewClockAttribute(file, chroma string) *ClockAttribute {
    clockAttr := &ClockAttribute{
        fileName: file,
        chromaName: chroma,
    }

    return clockAttr
}

func (clockAttr *ClockAttribute) String() string {
    return fmt.Sprintf("%s=%s#", clockAttr.chromaName, clockAttr.Value)
}

func (clockAttr *ClockAttribute) Encode() string {
    return fmt.Sprintf("%s %s;", clockAttr.fileName, clockAttr.Value)
}

func (clockAttr *ClockAttribute) Decode(s string) error {
    clockAttr.Value = strings.TrimPrefix(s, "string ")

    return nil
}

func (clockAttr *ClockAttribute) Update(edit Editor) error {
    var err error
    clockEdit, ok := edit.(*ClockEditor)
    if !ok {
        return fmt.Errorf("StringAttribute.Update requires StringEditor") 
    }

    clockAttr.Value, err = clockEdit.entry.GetText()
    return err
}

func (clock *ClockEditor) RunClock(cont func()) {
    state := PAUSE
    tick := time.NewTicker(time.Second)

    for {
        select {
        case state = <-clock.c:
        case <-tick.C:
        }    

        switch state {
        case START:
            // update time and animate
            currentText, err := clock.entry.GetText()
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
            clock.entry.SetText(currentTime.Format(clock.timeFormat))
            tick = time.NewTicker(time.Second)

            cont()
        case PAUSE:
            // block until we recieve an instruction
            state = <-clock.c
        case STOP:
            // reset the time and block
            clock.entry.SetText("00:00")
            state = <-clock.c
        default:
            log.Printf("Clock recieved unknown value through channel %d\n", state)
        }
    }
}

type ClockEditor struct {
    box    *gtk.Box
    entry  *gtk.Entry
    timeFormat   string
    c      chan int
}

func NewClockEditor(name string, animate, cont func()) (clockEdit *ClockEditor, err error) {
    clockEdit = &ClockEditor{timeFormat: "04:05"}

    clockEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        return
    }

    clockEdit.box.SetVisible(true)

    actions, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil { 
        return
    }

    clockEdit.box.PackStart(actions, false, false, padding)
    actions.SetVisible(true)

    startButton, err := gtk.ButtonNewWithLabel("Start")
    if err != nil { 
        return
    }

    actions.PackStart(startButton, false, false, padding)
    startButton.SetVisible(true)
    startButton.Connect("clicked", func() {
        clockEdit.c <- START
    })

    pauseButton, err := gtk.ButtonNewWithLabel("Pause")
    if err != nil { 
        return
    }

    actions.PackStart(pauseButton, false, false, padding)
    pauseButton.SetVisible(true)
    pauseButton.Connect("clicked", func() {
        clockEdit.c <- PAUSE 
    })

    stopButton, err := gtk.ButtonNewWithLabel("Stop")
    if err != nil { 
        return
    }

    actions.PackStart(stopButton, false, false, padding)
    stopButton.SetVisible(true)
    stopButton.Connect("clicked", func() {
        clockEdit.c <- STOP
    })

    go clockEdit.RunClock(cont)

    return
}

func (clockEdit *ClockEditor) Update(attr Attribute) error {
    clockAttr, ok := attr.(*ClockAttribute)
    if !ok {
        return fmt.Errorf("ClockEditor.Update requires ClockAttribute")
    }

    clockEdit.entry.SetText(clockAttr.Value)
    return nil
}

func (clockEdit *ClockEditor) Box() *gtk.Box {
    return clockEdit.box
}

