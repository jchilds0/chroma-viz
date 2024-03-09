package shows

import (
	"bufio"
	"chroma-viz/templates"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Show struct {
    numPages  int
    Pages     map[int]*Page
}

func NewShow() *Show {
    show := &Show{}
    show.Pages = make(map[int]*Page)
    return show
}

func (show *Show) SetPage(pageNum int, title string, temp *templates.Template, cont func(*Page)) {
    page := newPage(pageNum, title, temp, cont)

    if _, ok := show.Pages[pageNum]; ok {
        log.Printf("Page %d already exists", pageNum)
        return
    }

    show.Pages[pageNum] = page
}

func (show *Show) AddPage(title string, temp *templates.Template, cont func(*Page)) *Page {
    show.numPages++
    show.SetPage(show.numPages, title, temp, cont)

    return show.Pages[show.numPages]
}

func (show *Show) ImportShow(temps *templates.Temps, filename string, cont func(*Page)) error {
    pageReg, err := regexp.Compile("temp (?P<tempID>[0-9]*); title \"(?P<title>.*)\";")
    if err != nil {
        return err
    }

    propReg, err := regexp.Compile("index (?P<type>[0-9]*);")
    if err != nil {
        return err
    }

    file, err := os.Open(filename)
    if err != nil {
        return err
    }

    scanner := bufio.NewScanner(file)

    var page *Page
    for scanner.Scan() {
        line := scanner.Text()
        if pageReg.Match(scanner.Bytes()) {
            match := pageReg.FindStringSubmatch(line)
            tempID, err := strconv.Atoi(match[1])
            if err != nil {
                return err
            }

            temp := temps.Temps[tempID]
            show.AddPage(temp.Title, temp, cont)

            page = show.Pages[show.numPages]
            page.Title = match[2]
        } else if page != nil {
            match := propReg.FindStringSubmatch(line)
            if len(match) < 2 {
                log.Printf("Incorrect prop format (%s)\n", line)
                continue
            }

            index, err := strconv.Atoi(match[1])

            if err != nil {
                return err
            }

            prop := page.PropMap[index]

            if prop == nil {
                log.Printf("Unknown property (%d)\n", index)
                continue
            }

            prop.Decode(line)
        }
    }

    return nil
}

func (show *Show) ExportShow(filename string) {
    file, err := os.Create(filename)
    if err != nil {
        log.Fatalf("Error importing show (%s)", err)
    }
    defer file.Close()

    for _, page := range show.Pages {
        pageString := fmt.Sprintf("temp %d; title \"%s\";\n", page.TemplateID, page.Title)
        file.Write([]byte(pageString))

        for index, prop := range page.PropMap {
            if prop == nil {
                continue
            }

            file.WriteString(fmt.Sprintf("index %d;", index))

            file.WriteString(prop.Encode())

            file.WriteString("\n")
        }
    }
}
