package hub

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

func Logger(message string, args ...any) {
	file, err := os.OpenFile("./chroma_hub.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Error opening log file (%s)", err)
	}

	t := time.Now()
	s := fmt.Sprintf("[%s]\t", t.Format("2006-01-02 15:04:05")) + fmt.Sprintf(message, args...) + "\n"
	file.Write([]byte(s))
}

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

	for id, a := range archive.Assets {
		hub.Assets[id] = archive.Assets[id]

		log.Printf("Loaded Asset %d at %s/%s", id, a.Directory, a.Name)
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

func (hub *DataBase) ImportAssets(fileName string) (err error) {
	buf, err := os.ReadFile(fileName)
	if err != nil {
		return
	}

	var assets Assets
	err = json.Unmarshal(buf, &assets)

	for _, a := range assets {
		err = a.fetchImage()
		if err != nil {
			Logger("Error importing asset %d: %s", a.ImageID, err)
			continue
		}

		hub.Assets[a.ImageID] = a
	}

	return
}

func (hub *DataBase) ExportAssets(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		Logger("Couldn't open file: %s", err)
	}
	defer file.Close()

	buf, err := json.Marshal(hub.Assets)
	if err != nil {
		Logger("Error encoding assets: %s", err)
	}

	_, err = file.Write(buf)
	if err != nil {
		Logger("Error writing assets: %s", err)
	}

	log.Printf("Exported assets to %s", fileName)
}
