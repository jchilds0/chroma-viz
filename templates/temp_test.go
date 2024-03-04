package templates

import (
	"chroma-viz/props"
	"net"
	"testing"
	"time"

	"github.com/jchilds0/chroma-hub/chroma_hub"
)

func TestImportTemplates(t *testing.T) {
    temp := NewTemps()

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

    if len(temp.Temps) != 5 {
        t.Errorf("Incorrect num of templates, len(temp.Temps) = %d", len(temp.Temps))
    }

    for _, template := range temp.Temps {
        switch template.Title {
        case "Teal Box":
            propTest(t, template, 0, "Background", "rect")
            propTest(t, template, 1, "Circle", "circle")
            propTest(t, template, 2, "Title", "text")
            propTest(t, template, 3, "Subtitle", "text")
        case "Clock":
            propTest(t, template, 0, "Background", "rect")
            propTest(t, template, 1, "Circle", "circle")
            propTest(t, template, 2, "Left Split", "rect")
            propTest(t, template, 3, "Team 1", "text")
            propTest(t, template, 4, "Score 1", "text")
            propTest(t, template, 5, "Mid Split", "rect")
            propTest(t, template, 6, "Team 2", "text")
            propTest(t, template, 7, "Score 2", "text")
            propTest(t, template, 8, "Right Split", "rect")
            propTest(t, template, 9, "Clock", "clock")
            propTest(t, template, 10, "Period", "text")
        case "White Circle":
            propTest(t, template, 0, "Circle", "circle")
        case "Graph":
            propTest(t, template, 0, "Background", "rect")
            propTest(t, template, 1, "Graph", "graph")
            propTest(t, template, 2, "Title", "text")
        case "Ticker":
            propTest(t, template, 0, "Background", "rect")
            propTest(t, template, 1, "Box", "rect")
            propTest(t, template, 2, "Ticker", "ticker")
        default:
            t.Errorf("Unknown template %s", template.Title)
        }
    }
}

func propTest(t *testing.T, template *Template, i int, name, typed string) {
    prop := template.Prop[i]
    if prop.Name != name {
        t.Errorf("(%s) Incorrect prop name, expected %s, recieved %s", 
            template.Title, name, prop.Name)
    } 

    if prop.Type != props.StringToProp[typed] {
        t.Errorf("(%s) Incorrect prop type, expected %s, recieved %d", 
            template.Title, typed, prop.Type)
    }
}

func propTypeTest(t *testing.T, prop int, typed string) {
}
