package viz

import (
	"bufio"
	"chroma-viz/props"
	"chroma-viz/tcp"
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
    pages     map[int]*props.Page
    numPages  int
    edit      *Editor
    columns   [NUMCOL]bool
}

func NewShow(edit *Editor) *ShowTree {
    var err error
    show := &ShowTree{}

    show.TreeView, err = gtk.TreeViewNew()
    if err != nil {
        log.Fatalf("Error creating show (%s)", err)
    }

    show.TreeView.SetReorderable(true)

    show.pages = make(map[int]*props.Page)
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
    show.edit = edit

    // TODO: remove reference to show from outside scope
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

            show.edit.SetPage(show.pages[pageNum])
            conn["Preview"].SetPage <- show.pages[pageNum]
            conn["Preview"].SetAction <- tcp.ANIMATE_ON
        })

    return show 
}

func (show *ShowTree) NewShowPage(temp *props.Template) *props.Page {
    show.numPages++
    show.pages[show.numPages] = props.NewPage(show.numPages, temp.Title, temp)
    page := show.pages[show.numPages]
    show.treeList.Set(
        show.treeList.Append(), 
        []int{PAGENUM, TITLE, TEMPLATEID}, 
        []interface{}{page.PageNum, page.Title, page.TemplateID})
    return page
}

func (show *ShowTree) ImportShow(temp *TempTree, filename string) {
    pageReg, err := regexp.Compile("temp (?P<tempID>[0-9]*); title \"(?P<title>.*)\";")
    if err != nil {
        log.Fatalf("Error importing show (%s)", err)
    }

    propReg, err := regexp.Compile("index (?P<type>[0-9]*);")
    if err != nil {
        log.Fatalf("Error importing show (%s)", err)
    }

    file, err := os.Open(filename)
    if err != nil {
        log.Fatalf("Error importing show (%s)", err)
    }

    scanner := bufio.NewScanner(file)

    var page *props.Page
    for scanner.Scan() {
        line := scanner.Text()
        if pageReg.Match(scanner.Bytes()) {
            match := pageReg.FindStringSubmatch(line)
            tempID, err := strconv.Atoi(match[1])
            if err != nil {
                log.Fatalf("Error importing show (%s)", err)
            }

            page = show.NewShowPage(temp.temps[tempID])
            page.Title = match[2]
        } else if page != nil {
            match := propReg.FindStringSubmatch(line)
            if len(match) < 2 {
                log.Printf("Incorrect prop format (%s)\n", line)
                continue
            }

            index, err := strconv.Atoi(match[1])

            if err != nil {
                log.Fatalf("Error importing show (%s)", err);
            }

            prop := page.PropMap[index]

            if prop == nil {
                log.Printf("Unknown property (%d)\n", index)
                continue
            }

            prop.Decode(line)
        }
    }
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
