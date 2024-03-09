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
    c           chan int
    cont        func()
    currentTime string
    timeFormat  string
}

func NewClockAttribute(file, chroma string, cont func()) *ClockAttribute {
    clockAttr := &ClockAttribute{
        fileName: file,
        chromaName: chroma,
        cont: cont,
        timeFormat: "04:05",
        c: make(chan int),
    }

    go clockAttr.RunClock(cont)
    return clockAttr
}

func (clockAttr *ClockAttribute) String() string {
    return fmt.Sprintf("%s=%s#", clockAttr.chromaName, clockAttr.currentTime)
}

func (clockAttr *ClockAttribute) Encode() string {
    return fmt.Sprintf("%s %s;", clockAttr.fileName, clockAttr.currentTime)
}

func (clockAttr *ClockAttribute) Decode(s string) error {
    clockAttr.currentTime = strings.TrimPrefix(s, "string ")

    return nil
}

func (clockAttr *ClockAttribute) Update(edit Editor) error {
    var err error
    clockEdit, ok := edit.(*ClockEditor)
    if !ok {
        return fmt.Errorf("ClockAttribute.Update requires ClockEditor") 
    }

    clockAttr.currentTime, err = clockEdit.entry.GetText()
    return err
}

func (clock *ClockAttribute) RunClock(cont func()) {
    state := PAUSE
    tick := time.NewTicker(time.Second)

    for {
        select {
        case state = <-clock.c:
            tick = time.NewTicker(time.Second)
        case <-tick.C:
        }    

        switch state {
        case START:
            // update time and animate
            currentTime, err := time.Parse(clock.timeFormat, clock.currentTime)
            if err != nil {
                log.Println(err)
                continue
            }

            currentTime = currentTime.Add(time.Second)
            clock.currentTime = currentTime.Format(clock.timeFormat)

            cont()
        case PAUSE:
            // block until we recieve an instruction
            state = <-clock.c
        case STOP:
            // reset the time and block
            clock.currentTime = "00:00"
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

func NewClockEditor(name string, animate, cont func()) *ClockEditor {
    var err error
    clockEdit := &ClockEditor{timeFormat: "04:05"}

    clockEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        log.Print(err)
        return nil
    }

    clockEdit.box.SetVisible(true)

    actions, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil { 
        log.Print(err)
        return nil
    }

    clockEdit.box.PackStart(actions, false, false, padding)
    actions.SetVisible(true)

    startButton, err := gtk.ButtonNewWithLabel("Start")
    if err != nil { 
        log.Print(err)
        return nil
    }

    actions.PackStart(startButton, false, false, padding)
    startButton.SetVisible(true)
    startButton.Connect("clicked", func() {
        clockEdit.c <- START
    })

    pauseButton, err := gtk.ButtonNewWithLabel("Pause")
    if err != nil { 
        log.Print(err)
        return nil
    }

    actions.PackStart(pauseButton, false, false, padding)
    pauseButton.SetVisible(true)
    pauseButton.Connect("clicked", func() {
        clockEdit.c <- PAUSE 
    })

    stopButton, err := gtk.ButtonNewWithLabel("Stop")
    if err != nil { 
        log.Print(err)
        return nil
    }

    actions.PackStart(stopButton, false, false, padding)
    stopButton.SetVisible(true)
    stopButton.Connect("clicked", func() {
        clockEdit.c <- STOP
    })

    timeBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil {
        log.Print(err)
        return nil
    }

    timeBox.SetVisible(true)
    clockEdit.box.PackStart(timeBox, false, false, padding)

    label, err := gtk.LabelNew(name)
    if err != nil { 
        log.Print(err)
        return nil
    }

    label.SetVisible(true)
    label.SetWidthChars(12)
    timeBox.PackStart(label, false, false, padding)

    buf, err := gtk.EntryBufferNew("", 0)
    if err != nil { 
        log.Print(err)
        return nil
    }

    clockEdit.entry, err = gtk.EntryNewWithBuffer(buf)
    if err != nil { 
        log.Print(err)
        return nil
    }

    clockEdit.entry.SetVisible(true)
    clockEdit.entry.Connect("changed", animate)
    timeBox.PackStart(clockEdit.entry, false, false, 0)

    return clockEdit
}

func (clockEdit *ClockEditor) Update(attr Attribute) error {
    clockAttr, ok := attr.(*ClockAttribute)
    if !ok {
        return fmt.Errorf("ClockEditor.Update requires ClockAttribute")
    }

    clockEdit.c = clockAttr.c
    clockEdit.entry.SetText(clockAttr.currentTime)
    return nil
}

func (clockEdit *ClockEditor) Box() *gtk.Box {
    return clockEdit.box
}

