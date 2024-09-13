package hub

import (
	"chroma-viz/library/templates"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type DataBase struct {
	db        *sql.DB
	Templates map[int]*templates.Template
	Assets    map[int]Asset
}

type Templates struct {
	NumTemplates int
	Templates    []*templates.Template
}

func NewDataBase(numTemp int, username, password string) (hub *DataBase, err error) {
	hub = &DataBase{}
	hub.Assets = make(map[int]Asset, 128)
	hub.Templates = make(map[int]*templates.Template, 128)

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

func (hub *DataBase) GetTemplates() (temps []*templates.Template, err error) {
	rows, err := hub.db.Query("SELECT t.templateID FROM template t;")
	if err != nil {
		return
	}

	temps = make([]*templates.Template, 0, 10)

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
		temps = append(temps, temp)
	}

	return
}

func (hub *DataBase) CleanDB() {
	_, err := hub.db.Exec("DELETE FROM template")
	if err != nil {
		Logger("Error clearing db: %s", err)
	}

	return
}

func (hub *DataBase) TempIDs() (ids map[int]string, err error) {
	q := `
        SELECT t.templateID, t.Name
        FROM template t;
    `

	rows, err := hub.db.Query(q)
	if err != nil {
		return
	}

	ids = make(map[int]string)

	var (
		tempID int
		title  string
	)
	for rows.Next() {
		err = rows.Scan(&tempID, &title)
		if err != nil {
			return
		}

		ids[tempID] = title
	}

	return
}
