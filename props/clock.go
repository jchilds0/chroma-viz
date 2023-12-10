package props

import (
	"fmt"
	"log"
	"time"

	"github.com/gotk3/gotk3/gtk"
)


type ClockProp struct {
    box             *gtk.Box
    engine          bool
    preview         bool
    c               chan int
    editTime        *time.Time
    currentTime     time.Time
    timeFormat      string
    timeString      TextProp
}

func NewClockProp(width, height int, animate, cont func()) *ClockProp {
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

    clock.timeString = *NewTextProp(0, width, height, animate)
    clock.box.PackStart(clock.timeString.Tab(), false, false, 0)

    go clock.RunClock(cont)
    return clock
}

func (clock *ClockProp) Tab() *gtk.Box {
    return clock.box
}

func (clock *ClockProp) String() string {
    return fmt.Sprintf("text=%d#string=%s#pos_x=%d#pos_y=%d#", 
        clock.timeString.num, 
        clock.currentTime.Format(clock.timeFormat),
        clock.timeString.x_spin.GetValueAsInt(), 
        clock.timeString.y_spin.GetValueAsInt(),
    )
}

func (clock *ClockProp) Encode() string {
    return clock.timeString.Encode()
}

func (clock *ClockProp) Decode(input string) {
    clock.timeString.Decode(input)
    current, _ := clock.timeString.entry.GetText()
    edit, _ := time.Parse(clock.timeFormat, current)
    clock.editTime = &edit
}

func (clock *ClockProp) RunClock(cont func()) {
    state := PAUSE
    clock.currentTime = *clock.editTime

    for {
        select {
        case state = <-clock.c:
        default:
        }    

        switch state {
        case START:
            cont()
            clock.currentTime = clock.currentTime.Add(time.Second)
            time.Sleep(1 * time.Second)
        case PAUSE:
            // block until we recieve an instruction
            state = <-clock.c
        case STOP:
            // reset the time and block
            clock.currentTime = *clock.editTime
            cont()
            state = <-clock.c
        default:
            log.Printf("Clock recieved unknown value through channel %d\n", state)
        }
    }
}
