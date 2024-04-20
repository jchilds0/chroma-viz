package hub

import (
	"bufio"
	"chroma-viz/library/templates"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type DataBase struct {
	db        *sql.DB
	Templates map[int]*templates.Template
	Assets    map[int][]byte
	Dirs      map[int]string
	Names     map[int]string
}

func NewDataBase(numTemp int) *DataBase {
	hub := &DataBase{}
	hub.Templates = make(map[int]*templates.Template, numTemp)
	hub.Assets = make(map[int][]byte, 10)
	hub.Dirs = make(map[int]string, 10)
	hub.Names = make(map[int]string, 10)

	var err error
	hub.db, err = sql.Open("mysql", "/chroma_hub")
	if err != nil {
		log.Println(err)
	}

	return hub
}

// S -> {'num_temp': num, 'templates': [T]}
func (hub *DataBase) EncodeDB() (s string, err error) {
	var b strings.Builder

	first := true
	maxTempID := 0
	for _, temp := range hub.Templates {
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

func (hub *DataBase) TempIDs() (s string, err error) {
	q := `
        SELECT t.templateID, t.Name
        FROM template t;
    `

	rows, err := hub.db.Query(q)
	if err != nil {
		return
	}

	var (
		tempID int
		title  string
	)

	var b strings.Builder
	for rows.Next() {
		err = rows.Scan(&tempID, &title)
		if err != nil {
			return
		}

		b.WriteString(strconv.Itoa(tempID))
		b.WriteByte(' ')
		b.WriteString(title)
		b.WriteByte(';')
	}

	b.WriteString("EOF;")
	s = b.String()
	return
}

func (hub *DataBase) AddTemplate(id int64, name string, layer int) (err error) {
	// TODO: check for existing templates
	q := `
        INSERT INTO template VALUES (?, ?, ?);
    `

	_, err = hub.db.Exec(q, id, name, layer)
	return
}

func (hub *DataBase) AddGeometry(temp_id int64, name, geo_type string) (geo_id int64, err error) {
	q := `
        INSERT INTO geometry VALUES (NULL, ?, ?, ?);
    `

	result, err := hub.db.Exec(q, temp_id, name, geo_type)
	if err != nil {
		return
	}

	geo_id, err = result.LastInsertId()
	return
}

func (hub *DataBase) AddAttribute(geo_id int64, name, value, typed string, visible bool) (attr_id int64, err error) {
	q := `
        INSERT INTO attribute VALUES (NULL, ?, ?, ?, ?, ?);
    `

	result, err := hub.db.Exec(q, geo_id, name, value, typed, visible)
	if err != nil {
		return
	}

	attr_id, err = result.LastInsertId()
	return
}

func (hub *DataBase) GetTemplate(tempID int) (temp *templates.Template, err error) {
	tempQuery := `
        SELECT t.Name, t.Layer, COUNT(*)
        FROM template t
        INNER JOIN geometry g 
        ON g.templateID = t.templateID
        WHERE t.templateID = ?;
    `
	var (
		name    string
		layer   int
		num_geo int
	)

	row := hub.db.QueryRow(tempQuery, tempID)
	if err = row.Scan(&name, &layer, &num_geo); err != nil {
		log.Print(err)
		return
	}

	temp = templates.NewTemplate(name, tempID, layer, num_geo, 0)
	err = hub.GetGeometry(temp)
	return
}

func (hub *DataBase) GetGeometry(temp *templates.Template) (err error) {
	return
}

func (hub *DataBase) AcceptHubConn(ln net.Listener) {
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection (%s)", err)
			continue
		}

		go hub.HandleConn(conn)
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
func (hub *DataBase) HandleConn(conn net.Conn) {
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
			s, err = hub.EncodeDB()
			if err != nil {
				log.Print(err)
				continue
			}

			_, err = conn.Write([]byte(s))
		case "tempids":
			s, err = hub.TempIDs()
			if err != nil {
				log.Print(err)
				continue
			}

			_, err = conn.Write([]byte(s))
		case "temp":
			tempid, err := strconv.Atoi(cmds[4])
			if err != nil {
				break
			}

			template, err := hub.GetTemplate(tempid)
			if err == nil {
				log.Printf("Error getting template %d (%s)", tempid, err)
				continue
			}

			s, _ := template.Encode()
			_, err = conn.Write([]byte(s))
		case "img":
			imageID, _ := strconv.Atoi(cmds[4])
			image := hub.Assets[imageID]
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
			assets.Dirs = hub.Dirs
			assets.Names = hub.Names

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
