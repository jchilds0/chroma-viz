package viz

import (
	"bufio"
	"chroma-viz/shows"
	"chroma-viz/tcp"
	"chroma-viz/templates"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const NUMCOL = 3

const (
    PAGENUM = iota
    TITLE
    TEMPLATEID 
)

var KEYTITLE = map[int]string{
    PAGENUM: "Page Num",
    TITLE: "Title",
    TEMPLATEID: "Template ID",
}
type ShowTree struct {
    *gtk.TreeView
    treeList  *gtk.ListStore
    pages     map[int]*shows.Page
    numPages  int
    columns   [NUMCOL]bool
}

func NewShow(pageToEditor func(*shows.Page)) *ShowTree {
    var err error
    show := &ShowTree{}

    show.TreeView, err = gtk.TreeViewNew()
    if err != nil {
        log.Fatalf("Error creating show (%s)", err)
    }

    show.TreeView.SetReorderable(true)

    show.pages = make(map[int]*shows.Page)
    show.columns = [NUMCOL]bool{true, true, true}

    /* Columns */
    var column *gtk.TreeViewColumn
    for key := range show.columns {
        if show.columns[key] == false {
            continue
        }

        switch key {
        case TITLE:
            title, err := gtk.CellRendererTextNew()
            if err != nil {
                log.Fatalf("Error creating show (%s)", err)
            }

            title.SetProperty("editable", true)
            title.Connect("edited", 
                func(cell *gtk.CellRendererText, path string, text string) {
                    iter, err := show.treeList.GetIterFromString(path)
                    if err != nil {
                        log.Fatalf("Error editing page (%s)", err)
                    }

                    id, err := show.treeList.GetValue(iter, PAGENUM)
                    if err != nil {
                        log.Fatalf("Error editing page (%s)", err)
                    }

                    val, err := id.GoValue()
                    if err != nil {
                        log.Fatalf("Error editing page (%s)", err)
                    }

                    pageNum, err := strconv.Atoi(val.(string))
                    if err != nil {
                        log.Fatalf("Error editing page (%s)", err)
                    }

                    show.pages[pageNum].Title = text
                    show.treeList.SetValue(iter, TITLE, text)
            })

            column, err = gtk.TreeViewColumnNewWithAttribute(KEYTITLE[key], title, "text", 1)
            if err != nil {
                log.Fatalf("Error creating show (%s)", err)
            }

        default:
            cell, err := gtk.CellRendererTextNew()
            if err != nil {
                log.Fatalf("Error creating show (%s)", err)
            }

            column, err = gtk.TreeViewColumnNewWithAttribute(KEYTITLE[key], cell, "text", key)
            if err != nil {
                log.Fatalf("Error creating show (%s)", err)
            }

        }

        show.AppendColumn(column)
    }

    show.treeList, err = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
    if err != nil {
        log.Fatalf("Error creating show (%s)", err)
    }

    show.SetModel(show.treeList)

    show.Connect("row-activated", 
        func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) { 
            iter, err := show.treeList.GetIter(path)
            if err != nil {
                log.Fatalf("Error sending page to editor (%s)", err)
            }

            id, err := show.treeList.GetValue(iter, PAGENUM)
            if err != nil {
                log.Fatalf("Error sending page to editor (%s)", err)
            }

            val, err := id.GoValue()
            if err != nil {
                log.Fatalf("Error sending page to editor (%s)", err)
            }

            pageNum, err := strconv.Atoi(val.(string))
            if err != nil {
                log.Fatalf("Error sending page to editor (%s)", err)
            }

            pageToEditor(show.pages[pageNum])
            SendPreview(show.pages[pageNum], tcp.ANIMATE_ON)
        })

    return show 
}

func (show *ShowTree) NewShowPage(temp *templates.Template) *shows.Page {
    if temp == nil {
        log.Printf("Invalid template")
        return nil
    }

    show.numPages++
    show.pages[show.numPages] = shows.NewPage(show.numPages, temp.Title, temp)
    page := show.pages[show.numPages]
    show.treeList.Set(
        show.treeList.Append(), 
        []int{PAGENUM, TITLE, TEMPLATEID}, 
        []interface{}{page.PageNum, page.Title, page.TemplateID})
    return page
}

func (show *ShowTree) ImportPage(title string, temp *templates.Template) error {
    if temp == nil {
        return fmt.Errorf("Missing template")
    }

    show.numPages++
    show.pages[show.numPages] = shows.NewPage(show.numPages, title, temp)
    show.treeList.Set(
        show.treeList.Append(), 
        []int{PAGENUM, TITLE, TEMPLATEID}, 
        []interface{}{show.numPages, title, temp.TempID})
    return nil 
}

func (show *ShowTree) ImportShow(temp *TempTree, filename string) error {
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

    var page *shows.Page
    for scanner.Scan() {
        line := scanner.Text()
        if pageReg.Match(scanner.Bytes()) {
            match := pageReg.FindStringSubmatch(line)
            tempID, err := strconv.Atoi(match[1])
            if err != nil {
                return err
            }

            page = show.NewShowPage(temp.Temps[tempID])
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

func (show *ShowTree) ExportShow(filename string) {
    file, err := os.Create(filename)
    if err != nil {
        log.Fatalf("Error importing show (%s)", err)
    }
    defer file.Close()

    for _, page := range show.pages {
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
