package templates

import (
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
}
