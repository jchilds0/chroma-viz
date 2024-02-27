package shows

import (
	"chroma-viz/props"
	"chroma-viz/templates"
	"net"
	"testing"
	"time"

	"github.com/jchilds0/chroma-hub/chroma_hub"
)

func TestImportShow(t *testing.T) {
    fileName := "testing.show"
    temp := templates.NewTemps()
    show := NewShow()

    go chroma_hub.StartHub(9000, 2, "test.json")

    time.Sleep(1 * time.Second)
    conn, err := net.Dial("tcp", "127.0.0.1:9000")
    if err != nil {
        t.Fatalf("Error connecting to graphics hub (%s)", err)
    }

    err = temp.ImportTemplates(conn)
    if err != nil {
        t.Fatalf("Error importing graphics hub (%s)", err)
    }

    show.ImportShow(temp, fileName)

    if len(show.Pages) != 4 {
        t.Errorf("Incorrect number of pages (len(show.Pages) = %d)", len(show.Pages))
    }

    for _, page := range show.Pages {
        switch page.Title {
        case "Blue Box":
            rectPropTest(t, page.PropMap[0], 50, 100, 850, 180)
            circlePropTest(t, page.PropMap[1], 99, 90, 30, 75, 45, 315)
            textPropTest(t, page.PropMap[2], 190, 100, "Lower Frame")
            textPropTest(t, page.PropMap[3], 190, 30, "A lower frame subtitle")
        case "Clock Box":
            rectPropTest(t, page.PropMap[0], 50, 965, 810, 65)
            circlePropTest(t, page.PropMap[1], 35, 32, 0, 25, 0, 360)
            rectPropTest(t, page.PropMap[2], 76, 0, 10, 65)
            textPropTest(t, page.PropMap[3], 20, 15, "TEAM")
            textPropTest(t, page.PropMap[4], 180, 15, "0")
            rectPropTest(t, page.PropMap[5], 230, 0, 10, 65)
            textPropTest(t, page.PropMap[6], 20, 15, "TEAM")
            textPropTest(t, page.PropMap[7], 180, 15, "0")
            rectPropTest(t, page.PropMap[8], 230, 0, 10, 65)
            clockPropTest(t, page.PropMap[9], 110, 15)
            textPropTest(t, page.PropMap[10], 30, 15, "Q1")
        case "Ticker":
            rectPropTest(t, page.PropMap[0], 0, 0, 1920, 75)
            rectPropTest(t, page.PropMap[1], 1700, 25, 400, 100)
            tickerPropTest(t, page.PropMap[2], 25, 20, "Hello there", "world")
        case "Graph":
            rectPropTest(t, page.PropMap[0], 50, 100, 570, 240)
            graphPropTest(t, page.PropMap[1], 25, 25)
            textPropTest(t, page.PropMap[2], 50, 175, "Step Graph")
        default:
            t.Errorf("Unknown page %s", page.Title)
        }
    }

    t.Fatalf("Testing finished")
}

func rectPropTest(t *testing.T, prop props.Property, x, y, w, h int) {
    rect, ok := prop.(*props.RectProp)
    if !ok {
        t.Errorf("Prop %s is not a rect prop", prop.Name())
        return
    }

    if rect.Value[0] != x {
        t.Errorf("Rect prop incorrect value (rect.x = %d), expected %d", rect.Value[0], x)
    } else if rect.Value[1] != y {
        t.Errorf("Rect prop incorrect value (rect.y = %d), expected %d", rect.Value[1], y)
    } else if rect.Value[2] != w {
        t.Errorf("Rect prop incorrect value (rect.width = %d), expected %d", rect.Value[2], w)
    } else if rect.Value[3] != h {
        t.Errorf("Rect prop incorrect value (rect.height = %d), expected %d", rect.Value[3], h)
    }
}

func textPropTest(t *testing.T, prop props.Property, x, y int, text string) {
    textProp, ok := prop.(*props.TextProp)
    if !ok {
        t.Errorf("Prop %s is not a text prop", prop.Name())
        return
    }

    if textProp.Value[0] != x {
        t.Errorf("Text prop incorrect value (text.x = %d), expected %d", textProp.Value[0], x)
    } else if textProp.Value[1] != y {
        t.Errorf("Text prop incorrect value (text.y = %d), expected %d", textProp.Value[1], y)
    } else if textProp.S != text {
        t.Errorf("Text prop incorrect value (text.s = %s), expected %s", textProp.S, text)
    }
}

func circlePropTest(t *testing.T, prop props.Property, x, y, ir, or, sa, ea int) {
    circ, ok := prop.(*props.CircleProp)
    if !ok {
        t.Errorf("Prop %s is not a circle prop", prop.Name())
        return
    }

    if circ.Value[0] != x {
        t.Errorf("Circle prop incorrect value (circle.x = %d), expected %d", circ.Value[0], x)
    } else if circ.Value[1] != y {
        t.Errorf("Circle prop incorrect value (circle.y = %d), expected %d", circ.Value[1], y)
    } else if circ.Value[2] != ir {
        t.Errorf("Circle prop incorrect value (circle.inner_radius = %d), expected %d", circ.Value[2], ir)
    } else if circ.Value[3] != or {
        t.Errorf("Circle prop incorrect value (circle.outer_radius = %d), expected %d", circ.Value[3], or)
    } else if circ.Value[4] != sa {
        t.Errorf("Circle prop incorrect value (circle.start_angle = %d), expected %d", circ.Value[4], sa)
    } else if circ.Value[5] != ea {
        t.Errorf("Circle prop incorrect value (circle.end_angle = %d), expected %d", circ.Value[5], ea)
    }
}

func clockPropTest(t *testing.T, prop props.Property, x, y int) {
    clock, ok := prop.(*props.ClockProp)
    if !ok {
        t.Errorf("Prop %s is not a clock prop", prop.Name())
        return
    }

    if clock.Value[0] != x {
        t.Errorf("Clock prop incorrect value (clock.x = %d), expected %d", clock.Value[0], x)
    } else if clock.Value[1] != y {
        t.Errorf("Clock prop incorrect value (clock.y = %d), expected %d", clock.Value[1], y)
    }
}

func tickerPropTest(t *testing.T, prop props.Property, x, y int, s ...string) {
    tick, ok := prop.(*props.TickerProp)
    if !ok {
        t.Errorf("Prop %s is not a ticker prop", prop.Name())
        return
    }

    if tick.Value[0] != x {
        t.Errorf("Ticker prop incorrect value (tick.x = %d), expected %d", tick.Value[0], x)
    } else if tick.Value[1] != y {
        t.Errorf("Ticker prop incorrect value (tick.y = %d), expected %d", tick.Value[1], y)
    }

    // check list store values
}

func graphPropTest(t *testing.T, prop props.Property, x, y int) {
    g, ok := prop.(*props.GraphProp)
    if !ok {
        t.Errorf("Prop %s is not a graph prop", prop.Name())
        return
    }

    if g.Value[0] != x {
        t.Errorf("Graph prop incorrect value (g.x = %d), expected %d", g.Value[0], x)
    } else if g.Value[1] != y {
        t.Errorf("Graph prop incorrect value (g.y = %d), expected %d", g.Value[1], y)
    }

    // check list store values
}
