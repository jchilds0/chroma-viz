package attribute

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

const (
    START = iota
    PAUSE
    STOP
)

type ClockAttribute struct {
    Name        string 
    Type        int
    c           chan int
    CurrentTime string
    TimeFormat  string
}

func NewClockAttribute(name string, cont func()) *ClockAttribute {
    clockAttr := &ClockAttribute{
        Name: name,
        Type: CLOCK,
        TimeFormat: "04:05",
        c: make(chan int),
    }

    go clockAttr.RunClock(cont)
    return clockAttr
}

func (clockAttr *ClockAttribute) UnmarshalJSON(b []byte) error {
    var clockAttrJSON struct {
        ClockAttribute
        UnmarshalJSON struct {}
    }

    err := json.Unmarshal(b, &clockAttrJSON)
    if err != nil {
        return err
    }

    *clockAttr = clockAttrJSON.ClockAttribute
    clockAttr.c = make(chan int)
    return nil
}

func (clockAttr *ClockAttribute) SetClock(cont func()) {
    go clockAttr.RunClock(cont)
}

func (clockAttr *ClockAttribute) String() string {
    return fmt.Sprintf("%s=%s#", clockAttr.Name, clockAttr.CurrentTime)
}

func (clockAttr *ClockAttribute) Encode() string {
    return fmt.Sprintf("{'name': '%s', 'value': '%s'}", 
        clockAttr.Name, clockAttr.CurrentTime)
}

func (clockAttr *ClockAttribute) Update(edit Editor) error {
    var err error
    clockEdit, ok := edit.(*ClockEditor)
    if !ok {
        return fmt.Errorf("ClockAttribute.Update requires ClockEditor") 
    }

    clockAttr.CurrentTime, err = clockEdit.entry.GetText()
    return err
}

func (clock *ClockAttribute) RunClock(cont func()) {
    state := PAUSE
    tick := time.NewTicker(time.Second)

    for {
        select {
        case state = <-clock.c:
            tick.Reset(time.Second)
        case <-tick.C:
        }    

        switch state {
        case START:
            // update time and animate
            currentTime, err := time.Parse(clock.TimeFormat, clock.CurrentTime)
            if err != nil {
                log.Println(err)
                continue
            }

            currentTime = currentTime.Add(time.Second)
            clock.CurrentTime = currentTime.Format(clock.TimeFormat)
            cont()
        case PAUSE:
            // block until we recieve an instruction
            state = <-clock.c
        case STOP:
            // reset the time and block
            clock.CurrentTime = "00:00"
            state = <-clock.c
        default:
            log.Printf("Clock recieved unknown value through channel %d\n", state)
        }
    }
}

type ClockEditor struct {
    box    *gtk.Box
    entry  *gtk.Entry
    c      chan int
    name   string 
}

func NewClockEditor(name string) *ClockEditor {
    var err error
    clockEdit := &ClockEditor{name: name}

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
    timeBox.PackStart(clockEdit.entry, false, false, 0)

    return clockEdit
}

func (clockEdit *ClockEditor) Name() string {
    return clockEdit.name
}

func (clockEdit *ClockEditor) Update(attr Attribute) error {
    clockAttr, ok := attr.(*ClockAttribute)
    if !ok {
        return fmt.Errorf("ClockEditor.Update requires ClockAttribute")
    }

    clockEdit.c = clockAttr.c
    clockEdit.entry.SetText(clockAttr.CurrentTime)
    return nil
}

func (clockEdit *ClockEditor) Box() *gtk.Box {
    return clockEdit.box
}

