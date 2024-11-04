package hub

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

func Logger(message string, args ...any) {
	file, err := os.OpenFile("./log/chroma_hub.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Error opening log file (%s)", err)
	}

	t := time.Now()
	s := fmt.Sprintf("[%s]\t", t.Format("2006-01-02 15:04:05")) + fmt.Sprintf(message, args...) + "\n"
	file.Write([]byte(s))
}

func (hub *DataBase) ImportTemplates(fileName string) error {
	buf, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	var temp Templates
	err = json.Unmarshal(buf, &temp)
	if err != nil {
		return err
	}

	for _, temp := range temp.Templates {
		hub.ImportTemplate(temp)
		log.Printf("Loaded Template %d (%s)", temp.TempID, temp.Title)
	}

	//log.Printf("Imported Hub")
	return nil
}

func (hub *DataBase) ExportTemplates(fileName string) {
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

		hub.assets[a.ImageID] = a
	}

	return
}

func (hub *DataBase) ExportAssets(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		Logger("Couldn't open file: %s", err)
	}
	defer file.Close()

	buf, err := json.Marshal(hub.assets)
	if err != nil {
		Logger("Error encoding assets: %s", err)
	}

	_, err = file.Write(buf)
	if err != nil {
		Logger("Error writing assets: %s", err)
	}

	log.Printf("Exported assets to %s", fileName)
}
