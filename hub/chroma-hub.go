package hub

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

func StartHub(port, count int, fileName string) {
    hub := NewDataBase()
    ln, err := net.Listen("tcp", ":" + strconv.Itoa(port))
    if err != nil {
        log.Fatalf("Error creating server (%s)", err)
    }
    defer ln.Close()

    err = ImportArchive(hub, fileName)
    if err != nil {
        log.Printf("Error importing archive (%s)", err)
    }

    upperLimit := (count != -1)
    count = 2 * count

    for {
        if upperLimit && count == 0 {
            break
        }

        conn, err := ln.Accept()
        if err != nil {
            log.Printf("Error accepting connection (%s)", err)
        }

        handleConnection(hub, conn)

        log.Printf("Sent Hub to %s", conn.RemoteAddr())
        if upperLimit {
            count--
        }
    }
}

/*
    Send graphics hub to client using the following grammar

    S -> {'num_temp': num, 'templates': [T]}
    T -> {'id': num, 'num_geo': num, 'geometry': [G]} | T, T 
    G -> {'id': num, 'type': string, 'attr': [A]} | G, G 
    A -> {'name': string, 'value': string} | A, A

 */

func handleConnection(hub *DataBase, conn net.Conn) {
    _, err := conn.Write([]byte(hub.EncodeDB()))
    if err != nil {
        log.Printf("Error sending hub (%s)", err)
    }
}

func ImportArchive(hub *DataBase, fileName string) error {
    buf, err := os.ReadFile(fileName)
    if err != nil {
        return err
    }

    err = json.Unmarshal(buf, hub)
    if err != nil {
        return err
    }

    for _, temp := range hub.Templates {
        log.Printf("Loaded Template %d (%s)", temp.TempID, temp.Title)
    }

    log.Printf("Imported Hub")
    return nil
}

func ExportArchive(hub *DataBase, fileName string) {
    file, err := os.Create(fileName)
    if err != nil {
        log.Fatalf("Couldn't open file (%s)", err)
    }
    defer file.Close()

    buf, err := json.Marshal(hub)
    if err != nil {
        log.Printf("Error encoding hub (%s)", err)
    }

    _, err = file.Write(buf)
    if err != nil {
        log.Printf("Error writing hub (%s)", err)
    }
    
    fmt.Printf("Exported Hub\n")
}

