package props

import (
	"fmt"
	"log"
	"time"

	"github.com/gotk3/gotk3/gtk"
)


type ClockProp struct {
    box             *gtk.Box
    name            string
    engine          bool
    preview         bool
    c               chan int
    editTime        *time.Time
    currentTime     time.Time
    timeFormat      string
    timeString      TextProp
}

func NewClockProp(width, height int, animate, cont func(), name string) *ClockProp {
    var err error
    clock := &ClockProp{name: name}
    clock.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    clock.box.SetVisible(true)
    clock.c = make(chan int, 1)

    actions, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    clock.box.PackStart(actions, false, false, padding)
    actions.SetVisible(true)

    startButton, err := gtk.ButtonNewWithLabel("Start")
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    actions.PackStart(startButton, false, false, padding)
    startButton.SetVisible(true)
    startButton.Connect("clicked", func() {
        clock.c <- START
    })

    pauseButton, err := gtk.ButtonNewWithLabel("Pause")
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    actions.PackStart(pauseButton, false, false, padding)
    pauseButton.SetVisible(true)
    pauseButton.Connect("clicked", func() {
        clock.c <- PAUSE 
    })

    stopButton, err := gtk.ButtonNewWithLabel("Stop")
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    actions.PackStart(stopButton, false, false, padding)
    stopButton.SetVisible(true)
    stopButton.Connect("clicked", func() {
        clock.c <- STOP
    })

    clock.timeFormat = "04:05"
    edit, err := time.Parse(clock.timeFormat, "00:00")
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    clock.editTime = &edit

    clock.timeString = *NewTextProp(width, height, animate, "clock time")
    clock.box.PackStart(clock.timeString.Tab(), false, false, 0)

    go clock.RunClock(cont)
    return clock
}

func (clock *ClockProp) Name() string {
    return clock.name
}

func (clock *ClockProp) Tab() *gtk.Box {
    return clock.box
}

func (clock *ClockProp) String() string {
    currentTime := clock.currentTime.Format(clock.timeFormat)
    currentString := fmt.Sprintf("string=%s#pos_x=%d#pos_y=%d#", 
        currentTime,
        clock.timeString.x_spin.GetValueAsInt(), 
        clock.timeString.y_spin.GetValueAsInt(),
    )
    
    return currentString
}

func (clock *ClockProp) Encode() string {
    return clock.timeString.Encode()
}

func (clock *ClockProp) Decode(input string) {
    clock.timeString.Decode(input)
    current, err := clock.timeString.entry.GetText()
    if err != nil { 
        log.Printf("Error decoding clock (%s)", err) 
    }

    edit, err := time.Parse(clock.timeFormat, current)
    if err != nil { 
        log.Printf("Error decoding clock (%s)", err) 
    }

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
