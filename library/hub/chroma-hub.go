package hub

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func Logger(message string, args ...any) {
	file, err := os.Create("./chroma_hub.log")
	if err != nil {
		log.Printf("Error opening log file (%s)", err)
	}

	t := time.Now()
	s := fmt.Sprintf("[%s]\t", t.Format("2006-01-02 15:04:05")) + fmt.Sprintf(message, args...)
	file.Write([]byte(s))
}

/*
Send graphics hub to client using the following grammar

S -> {'num_temp': num, 'templates': [T]}
T -> {'id': num, 'num_geo': num, 'geometry': [G]} | T, T
G -> {'id': num, 'type': string, 'attr': [A]} | G, G
A -> {'name': string, 'value': string} | A, A
*/
func (hub *DataBase) ImportArchive(fileName string) error {
	buf, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	var archive DataBase
	err = json.Unmarshal(buf, &archive)
	if err != nil {
		return err
	}

	for _, temp := range archive.Templates {
		if _, ok := hub.Templates[temp.TempID]; ok {
			return fmt.Errorf("Template ID %d already exists", temp.TempID)
		}

		hub.ImportTemplate(*temp)
		log.Printf("Loaded Template %d (%s)", temp.TempID, temp.Title)
	}

	for id := range archive.Assets {
		hub.Assets[id] = archive.Assets[id]
		hub.Dirs[id] = archive.Dirs[id]
		hub.Names[id] = archive.Names[id]

		log.Printf("Loaded Asset %d at %s/%s", id, archive.Dirs[id], archive.Names[id])
	}

	//log.Printf("Imported Hub")
	return nil
}

func (hub *DataBase) ExportArchive(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		Logger("Couldn't open file (%s)", err)
	}
	defer file.Close()

	buf, err := json.Marshal(hub)
	if err != nil {
		Logger("Error encoding hub (%s)", err)
	}

	_, err = file.Write(buf)
	if err != nil {
		Logger("Error writing hub (%s)", err)
	}

	log.Printf("Exported hub to %s", fileName)
}

func (hub *DataBase) ExportTemplate(fileName string, tempID int64) error {
	temp, err := hub.GetTemplate(tempID)
	if err != nil {
		return err
	}

	err = temp.ExportTemplate(fileName)
	return err
}

func (hub *DataBase) ImportAsset(args []string) (err error) {
	if len(args) != 4 {
		return fmt.Errorf("Incorrect number of args to import asset")
	}

	image, err := os.ReadFile(args[0])
	if err != nil {
		return
	}

	imageID, err := strconv.Atoi(args[3])
	if err != nil {
		return
	}

	hub.Assets[imageID] = image
	hub.Dirs[imageID] = args[1]
	hub.Names[imageID] = args[2]
	return
}
