package hub

import (
	"chroma-viz/library/props"
	"chroma-viz/library/templates"
	"fmt"
	"log"
	"net"
)

type DataBase struct {
    Templates map[int]*templates.Template
}

func NewDataBase() *DataBase {
    db := &DataBase{}
    db.Templates = make(map[int]*templates.Template)

    return db
}

// S -> {'num_temp': num, 'templates': [T]}
func (db *DataBase) EncodeDB() string {
    first := true 
    templates := ""
    for _, temp := range db.Templates {
        if first {
            templates = temp.Encode()
            first = false 
            continue
        }

        templates = fmt.Sprintf("%s,%s", templates, temp.Encode())
    }

    return fmt.Sprintf("{'num_temp': %d, 'templates': [%s]}", len(db.Templates), templates)
}

func (db *DataBase) AddTemplate(id int, anim_on, anim_cont, anim_off string) {
    if db.Templates[id] != nil {
        log.Printf("Template %d already exists", id)
        return
    }

    db.Templates[id] = templates.NewTemplate("", id, 0, 10)
}

func (db *DataBase) AddGeometry(temp_id, geo_id int, geo_type string) {
    if db.Templates[temp_id] == nil {
        log.Printf("Template %d does not exist", temp_id)
    }

    temp := db.Templates[temp_id]
    temp.AddProp("", geo_id, props.StringToProp[geo_type], nil)
}

func (db *DataBase) SendHub(ln net.Listener) {
    // upperLimit := (count != -1)
    // count = 2 * count
    defer ln.Close()

    for {
        // if upperLimit && count == 0 {
        //     break
        // }

        conn, err := ln.Accept()
        if err != nil {
            log.Printf("Error accepting connection (%s)", err)
        }

        _, err = conn.Write([]byte(hub.EncodeDB()))
        if err != nil {
            log.Printf("Error sending hub (%s)", err)
        }

        log.Printf("Sent Hub to %s", conn.RemoteAddr())
        // if upperLimit {
        //     count--
        // }
    }
}