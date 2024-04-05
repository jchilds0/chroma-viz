package hub

import (
	"bufio"
	"chroma-viz/library/props"
	"chroma-viz/library/templates"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type DataBase struct {
	Templates map[int]*templates.Template
	Assets    map[int][]byte
	Dirs      map[int]string
	Names     map[int]string
}

func NewDataBase(numTemp int) *DataBase {
	db := &DataBase{}
	db.Templates = make(map[int]*templates.Template, numTemp)
	db.Assets = make(map[int][]byte, 10)
	db.Dirs = make(map[int]string, 10)
	db.Names = make(map[int]string, 10)

	return db
}

// S -> {'num_temp': num, 'templates': [T]}
func (db *DataBase) EncodeDB() (s string, err error) {
	var b strings.Builder

	first := true
	maxTempID := 0
	for _, temp := range db.Templates {
		maxTempID = max(maxTempID, temp.TempID)

		if !first {
			b.WriteString(",")
		}

		first = false
		tempStr, _ := temp.Encode()
		b.WriteString(tempStr)
	}

	s = fmt.Sprintf("{'num_temp': %d, 'templates': [%s]}", maxTempID+2, b.String())
	return
}

func (db *DataBase) TempIDs() (s string) {
	for _, temp := range db.Templates {
		if temp == nil {
			continue
		}

		s += fmt.Sprintf("%d %s;", temp.TempID, temp.Title)
	}

	return s + "EOF;"
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
	temp.AddGeometry("", geo_id, props.StringToProp[geo_type], nil)
}

func (db *DataBase) AcceptHubConn(ln net.Listener) {
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection (%s)", err)
			continue
		}

		go db.HandleConn(conn)
	}
}

/*
Protocol for Chroma Hub <=> Chroma Viz/Engine communication

S -> V C;
V -> ver %d %d
C -> full | tempids | temp %d

A command consists of a header with the protocol version,
the command which is currently either full or a single page,
and a template id if the command is template.

	full - encode the entire chroma hub and send to client

	tempids - send all current template ids

	temp i - encode template with template id 'i' and send
	to clint
*/
func (db *DataBase) HandleConn(conn net.Conn) {
	req := make([]byte, 0, 1024)
	buf := bufio.NewReader(conn)

	for {
		s, err := buf.ReadString(';')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error reading request (%s)", err)
			continue
		}

		cmds := strings.Split(strings.TrimSuffix(s, ";"), " ")

		if len(s) < 4 {
			continue
		}

		if cmds[1] != "0" || cmds[2] != "1" {
			log.Fatalf("Request has incorrect ver %s %s, expected 0 1 (%s)", cmds[1], cmds[2], s)
			continue
		}

		switch cmds[3] {
		case "full":
			s, _ = db.EncodeDB()
			_, err = conn.Write([]byte(s))
		case "tempids":
			_, err = conn.Write([]byte(db.TempIDs()))
		case "temp":
			tempid, err := strconv.Atoi(cmds[4])
			if err != nil {
				break
			}

			template := db.Templates[tempid]
			if template == nil {
				log.Printf("Template %d does not exist", tempid)
				continue
			}

			s, _ := template.Encode()
			_, err = conn.Write([]byte(s))
		case "img":
			imageID, _ := strconv.Atoi(cmds[4])
			image := db.Assets[imageID]
			if image == nil {
				log.Printf("image %d does not exist", imageID)
				_, err = conn.Write([]byte{0, 1, 0, 0})
				_, err = conn.Write([]byte{0, 0, 0, 0})
				continue
			}

			lenByte0 := byte(len(image) & (1<<8 - 1))
			lenByte1 := byte((len(image) >> 8) & (1<<8 - 1))
			lenByte2 := byte((len(image) >> 16) & (1<<8 - 1))
			lenByte3 := byte((len(image) >> 24) & (1<<8 - 1))

			_, err = conn.Write([]byte{0, 1, 0, 0})
			_, err = conn.Write([]byte{lenByte3, lenByte2, lenByte1, lenByte0})
			_, err = conn.Write(image)
		case "assets":
			var assets struct {
				Dirs  map[int]string
				Names map[int]string
			}
			assets.Dirs = db.Dirs
			assets.Names = db.Names

			assetsJson, err := json.Marshal(assets)
			if err != nil {
				break
			}

			_, err = conn.Write(assetsJson)
			_, err = conn.Write([]byte{0})
		default:
			log.Printf("Unknown request %s", string(req[:]))
			continue
		}

		if err != nil {
			log.Printf("Error responding to request %s (%s)", string(req[:]), err)
			continue
		}
	}
}
