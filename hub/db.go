package hub

import (
	"bufio"
	"chroma-viz/library/props"
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

var geoTables = map[int]string{
	props.RECT_PROP:   "rectangle",
	props.TEXT_PROP:   "text",
	props.CIRCLE_PROP: "circle",
	props.IMAGE_PROP:  "asset",
}

type DataBase struct {
	db        *sql.DB
	Templates map[int64]*templates.Template
	Assets    map[int][]byte
	Dirs      map[int]string
	Names     map[int]string
}

func NewDataBase(numTemp int) *DataBase {
	hub := &DataBase{}
	hub.Templates = make(map[int64]*templates.Template, 100)
	hub.Assets = make(map[int][]byte, 10)
	hub.Dirs = make(map[int]string, 10)
	hub.Names = make(map[int]string, 10)

	var err error
	hub.db, err = sql.Open("mysql", "/chroma_hub")
	if err != nil {
		Logger("Error connecting to database (%s) ", err)
	}

	return hub
}

// S -> {'num_temp': num, 'templates': [T]}
func (hub *DataBase) EncodeDB() (s string, err error) {
	var b strings.Builder
	first := true
	rows, err := hub.db.Query("SELECT t.templateID FROM template t;")
	if err != nil {
		return
	}

	var (
		maxTempID int64
		tempID    int64
		temp      *templates.Template
		tempStr   string
	)
	for rows.Next() {
		err = rows.Scan(&tempID)
		if err != nil {
			err = fmt.Errorf("Scan TempID: %s", err)
			return
		}

		maxTempID = max(maxTempID, tempID)

		if !first {
			b.WriteString(",")
		}

		first = false
		temp, err = hub.GetTemplate(tempID)
		if err != nil {
			err = fmt.Errorf("Retrieve Template: %s", err)
			return
		}

		tempStr, err = temp.Encode()
		if err != nil {
			err = fmt.Errorf("Encode Template: %s", err)
			return
		}

		b.WriteString(tempStr)
	}

	s = fmt.Sprintf("{'num_temp': %d, 'templates': [%s]}", maxTempID+2, b.String())
	return
}

func (hub *DataBase) CleanDB() (err error) {
	_, err = hub.db.Exec("DELETE FROM template")
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

func (hub *DataBase) ImportTemplate(temp templates.Template) (err error) {
	err = hub.AddTemplate(temp.TempID, temp.Title, temp.Layer)
	if err != nil {
		return
	}

	for _, geo := range temp.Geometry {
		geom := geo.Geom()

		switch geom.GeoType {
		case templates.GEO_RECT:
			rect, ok := geo.(*templates.Rectangle)
			if !ok {
				continue
			}

			hub.AddRectangle(temp.TempID, *rect)

		case templates.GEO_CIRCLE:
			circle, ok := geo.(*templates.Circle)
			if !ok {
				continue
			}

			hub.AddCircle(temp.TempID, *circle)

		case templates.GEO_TEXT:
			text, ok := geo.(*templates.Text)
			if !ok {
				continue
			}

			hub.AddText(temp.TempID, *text)

		default:
			log.Printf("Unknown geometry type")
		}
	}

	return
}

func (hub *DataBase) AddTemplate(tempID int64, name string, layer int) (err error) {
	// TODO: run as a transaction
	deleteTemp := `
        DELETE FROM template WHERE templateID = ?;
    `
	_, err = hub.db.Exec(deleteTemp, tempID)
	if err != nil {
		Logger(err.Error())
	}

	addTemp := `
        INSERT INTO template VALUES (?, ?, ?);
    `

	_, err = hub.db.Exec(addTemp, tempID, name, layer)
	return
}

func (hub *DataBase) GetTemplate(tempID int64) (temp *templates.Template, err error) {
	temp, ok := hub.Templates[tempID]
	if ok {
		return
	}

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
		return
	}

	temp = templates.NewTemplate(name, tempID, layer, num_geo, 0)
	err = hub.GetGeometry(temp)

	hub.Templates[tempID] = temp
	return
}

func (hub *DataBase) StartHub(port int) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		Logger("Error creating server (%s)", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			Logger("Error accepting connection (%s)", err)
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
			Logger("Error reading request (%s)", err)
			continue
		}

		cmds := strings.Split(strings.TrimSuffix(s, ";"), " ")

		if len(s) < 4 {
			continue
		}

		if cmds[1] != "0" || cmds[2] != "1" {
			Logger("Request has incorrect ver %s %s, expected 0 1 (%s)", cmds[1], cmds[2], s)
			continue
		}

		switch cmds[3] {
		case "full":
			s, err = hub.EncodeDB()
			if err != nil {
				Logger("Error retrieving database (%s)", err)
				continue
			}

			_, err = conn.Write([]byte(s))
		case "tempids":
			s, err = hub.TempIDs()
			if err != nil {
				Logger("Error retrieving template IDs (%s)", err)
				continue
			}

			_, err = conn.Write([]byte(s + string(props.END_OF_MESSAGE)))
		case "temp":
			tempid, err := strconv.ParseInt(cmds[4], 10, 64)
			if err != nil {
				Logger("Error getting template id (%s)", err)
				continue
			}

			template, err := hub.GetTemplate(tempid)
			if err != nil {
				Logger("Error getting template %d (%s)", tempid, err)
				continue
			}

			s, _ := template.Encode()
			_, err = conn.Write([]byte(s + string(props.END_OF_MESSAGE)))
		case "img":
			imageID, _ := strconv.Atoi(cmds[4])
			image := hub.Assets[imageID]
			if image == nil {
				Logger("Image %d does not exist", imageID)
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
			Logger("Unknown request %s", string(req[:]))
			continue
		}

		if err != nil {
			Logger("Error responding to request %s (%s)", string(req[:]), err)
			continue
		}
	}
}
