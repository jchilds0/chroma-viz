package shows

import (
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

    time.Sleep(5 * time.Second)
    conn, err := net.Dial("tcp", "127.0.0.1:9000")
    if err != nil {
        t.Fatalf("Error connecting to graphics hub (%s)", err)
    }

    err = temp.ImportTemplates(conn)
    if err != nil {
        t.Fatalf("Error importing graphics hub (%s)", err)
    }

    show.ImportShow(temp, fileName)
}
