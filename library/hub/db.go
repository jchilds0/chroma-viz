package hub

import (
	"bufio"
	"chroma-viz/library/templates"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type DataBase struct {
	db        *sql.DB
	Templates map[int64]*templates.Template
	Assets    map[int][]byte
	Dirs      map[int]string
	Names     map[int]string
}

func NewDataBase(numTemp int, username, password string) (hub *DataBase, err error) {
	hub = &DataBase{}
	hub.Templates = make(map[int64]*templates.Template, 100)
	hub.Assets = make(map[int][]byte, 10)
	hub.Dirs = make(map[int]string, 10)
	hub.Names = make(map[int]string, 10)

	hub.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/?multiStatements=true", username, password))
	if err != nil {
		err = fmt.Errorf("Error opening database: %s", err)
		return
	}

	err = hub.db.Ping()
	if err != nil {
		err = fmt.Errorf("Error connecting to database: %s", err)
		return
	}

	return
}

func (hub *DataBase) ImportSchema(filename string) (err error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	s := string(buf)

	_, err = hub.db.Exec(s)
	return
}

func (hub *DataBase) SelectDatabase(name, username, password string) (err error) {
	hub.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", username, password, name))
	return
}

// S -> {'num_temp': num, 'templates': [T]}
func (hub *DataBase) EncodeDB() (buf []byte, err error) {
	rows, err := hub.db.Query("SELECT t.templateID FROM template t;")
	if err != nil {
		return
	}

	var databaseJSON struct {
		NumTemplates int
		Templates    []*templates.Template
	}

	databaseJSON.Templates = make([]*templates.Template, 0, 10)

	var (
		maxTempID int64
		tempID    int64
		temp      *templates.Template
	)
	for rows.Next() {
		err = rows.Scan(&tempID)
		if err != nil {
			err = fmt.Errorf("TempID: %s", err)
			return
		}

		temp, err = hub.GetTemplate(tempID)
		if err != nil {
			err = fmt.Errorf("Retrieve Template: %s", err)
			return
		}

		maxTempID = max(maxTempID, temp.TempID)
		databaseJSON.Templates = append(databaseJSON.Templates, temp)
	}

	databaseJSON.NumTemplates = int(maxTempID)
	return json.Marshal(databaseJSON)
}

func (hub *DataBase) CleanDB() {
	_, err := hub.db.Exec("DELETE FROM template")
	if err != nil {
		Logger("Error clearing db: %s", err)
	}

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

var EOM = 6

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
	var b []byte
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
			_, err = conn.Write([]byte(string(EOM)))
			continue
		}

		switch cmds[3] {
		case "full":
			b, err = hub.EncodeDB()
			if err != nil {
				Logger("Error retrieving database (%s)", err)
				_, err = conn.Write([]byte(string(EOM)))
				continue
			}

			_, err = conn.Write(append(b, byte(EOM)))

		case "tempids":
			s, err = hub.TempIDs()
			if err != nil {
				Logger("Error retrieving template IDs (%s)", err)
				_, err = conn.Write([]byte(string(EOM)))
				continue
			}

			_, err = conn.Write([]byte(s + string(EOM)))

		case "temp":
			tempid, err := strconv.ParseInt(cmds[4], 10, 64)
			if err != nil {
				Logger("Error getting template id: %s", err)
				_, err = conn.Write([]byte(string(EOM)))
				continue
			}

			template, err := hub.GetTemplate(tempid)
			if err != nil {
				Logger("Error getting template %d: %s", tempid, err)
				_, err = conn.Write([]byte(string(EOM)))
				continue
			}

			b, err = json.Marshal(template)
			if err != nil {
				Logger("Error encoding template %d: %s", tempid, err)
				_, err = conn.Write([]byte(string(EOM)))
				continue
			}

			_, err = conn.Write(append(b, byte(EOM)))

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
