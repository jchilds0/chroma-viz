package props

import (
	"fmt"
	"log"
	"time"

	"github.com/gotk3/gotk3/gtk"
)


type ClockProp struct {
    box *gtk.Box
    active bool
    c chan int
    currentTime time.Time
    editTime *time.Time
    timeFormat string
}

func NewClockProp(animate func()) *ClockProp {
    clock := &ClockProp{}
    clock.box, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    clock.box.SetVisible(true)
    clock.c = make(chan int, 1)

    actions, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    clock.box.PackStart(actions, false, false, padding)
    actions.SetVisible(true)

    startButton, _ := gtk.ButtonNewWithLabel("Start")
    actions.PackStart(startButton, false, false, padding)
    startButton.SetVisible(true)
    startButton.Connect("clicked", func() {
        clock.c <- START
    })

    pauseButton, _ := gtk.ButtonNewWithLabel("Pause")
    actions.PackStart(pauseButton, false, false, padding)
    pauseButton.SetVisible(true)
    pauseButton.Connect("clicked", func() {
        clock.c <- PAUSE 
    })

    stopButton, _ := gtk.ButtonNewWithLabel("Stop")
    actions.PackStart(stopButton, false, false, padding)
    stopButton.SetVisible(true)
    stopButton.Connect("clicked", func() {
        clock.c <- STOP
    })

    clock.timeFormat = "04:05"
    edit, _ := time.Parse(clock.timeFormat, "00:00")
    clock.editTime = &edit

    inputBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    inputBox.SetVisible(true)
    clock.box.PackStart(inputBox, false, false, padding)

    labelTime, _ := gtk.LabelNew("Start Time: ")
    labelTime.SetVisible(true)
    labelTime.SetWidthChars(7)
    inputBox.PackStart(labelTime, false, false, padding)

    editTime := clock.editTime.Format(clock.timeFormat)
    bufTime, _ := gtk.EntryBufferNew(editTime, len(editTime))
    textTime, _ := gtk.EntryNewWithBuffer(bufTime)
    textTime.SetVisible(true)
    inputBox.PackStart(textTime, false, false, 0)

    textTime.Connect("changed", 
        func(e *gtk.Entry) { 
            text, err := e.GetText()
            newTime, err := time.Parse(clock.timeFormat, text)
            if err != nil {
                log.Printf("Incorrect time format entered (%s)\n", text)
                return
            }
            *clock.editTime = newTime
            animate()
        })

    go clock.RunClock(animate)
    return clock
}

func (clock *ClockProp) Tab() *gtk.Box {
    return clock.box
}

func (clock *ClockProp) String() string {
    return fmt.Sprintf("text0#%s#", 
        clock.currentTime.Format(clock.timeFormat))
}

func (clock *ClockProp) Encode() string {
    return ""
}

func (clock *ClockProp) Decode(input string) {
}

func (clock *ClockProp) RunClock(animate func()) {
    state := PAUSE
    clock.currentTime = *clock.editTime

    for {
        select {
        case state = <-clock.c:
        default:
        }    

        switch state {
        case START:
            animate()
            clock.currentTime = clock.currentTime.Add(time.Second)
            time.Sleep(1 * time.Second)
        case PAUSE:
            // block until we recieve an instruction
            state = <-clock.c
        case STOP:
            // reset the time and block
            clock.currentTime = *clock.editTime
            animate()
            state = <-clock.c
        default:
            log.Printf("Clock recieved unknown value through channel %d\n", state)
        }
    }
}
