package hub

import (
	"fmt"
	"log"
)

type DataBase struct {
    Array map[int]*graphics.Template
}

func NewDataBase() *DataBase {
    db := &DataBase{}
    db.Array = make(map[int]*graphics.Template)

    return db
}

func (db *DataBase) String() string {
    first := true 
    templates := ""
    for _, temp := range db.Array {
        if first {
            templates = temp.String()
            first = false 
            continue
        }

        templates = fmt.Sprintf("%s,%s", templates, temp.String())
    }

    return fmt.Sprintf("{'num_temp': %d, 'templates': [%s]}", len(db.Array), templates)
}

func (db *DataBase) AddTemplate(id int, anim_on, anim_cont, anim_off string) {
    if db.Array[id] != nil {
        log.Printf("Template %d already exists", id)
        return
    }

    db.Array[id] = graphics.NewTemplate(id, anim_on, anim_cont, anim_off)
}

func (db *DataBase) AddGeometry(temp_id, geo_id int, geo_type string) {
    if db.Array[temp_id] == nil {
        log.Printf("Template %d does not exist", temp_id)
    }

    temp := db.Array[temp_id]
    temp.AddGeometry(geo_id, geo_type)
}

func (db *DataBase) AddAttr(temp_id, geo_id int, name, attr string) {
    if db.Array[temp_id] == nil {
        log.Printf("Template %d does not exist", temp_id)
        return
    }

    db.Array[temp_id].AddAttr(geo_id, name, attr)
}
