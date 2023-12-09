package gui

import (
	"bufio"
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

type Page struct {
    Box         *gtk.ListBoxRow
    pageNum     int
    title       string
    templateID  int
    props       map[string]Property
}

func NewPage(pageNum int, title string, temp *Template) *Page {
    page := &Page{pageNum: pageNum, title: title, templateID: temp.templateID}
    page.props = make(map[string]Property)

    animate := func() { conn["Preview"].SendPage(page, ANIMATE_ON) }


    num_text := 1
    for name, prop := range temp.props {
        switch (prop) {
        case "RectProp":
            page.props[name] = NewRectProp(1920, 1080, animate)
        case "TextProp":
            count, _ := strconv.Atoi(name)
            page.props["Text " + strconv.Itoa(num_text)] = NewTextProp(count, animate)
            num_text++
        case "ClockProp":
            page.props[name] = NewClockProp(page, animate)
        default:
            log.Printf("Page %d: Unknown property %s", pageNum, prop)
        }
    }

    return page
}

func (page *Page) pageToListRow() *gtk.ListBoxRow {
    row1, _ := gtk.ListBoxRowNew()
    row1.Add(textToBuffer(strconv.Itoa(page.pageNum)))
    row1.Add(textToBuffer(page.title))

    return row1
}

type ShowTree struct {
    *gtk.TreeView
    treeList  *gtk.ListStore
    pages     map[int]*Page
    numPages  int
    edit      *Editor
    columns   [NUMCOL]bool
}

func NewShow(edit *Editor) *ShowTree {
    show := &ShowTree{}

    show.TreeView, _ = gtk.TreeViewNew()
    show.TreeView.SetReorderable(true)

    show.pages = make(map[int]*Page)
    show.columns = [NUMCOL]bool{true, true, true}

    /* Columns */
    var column *gtk.TreeViewColumn
    for key := range show.columns {
        if show.columns[key] == false {
            continue
        }

        switch key {
        case TITLE:
            title, _ := gtk.CellRendererTextNew()
            title.SetProperty("editable", true)
            title.Connect("edited", 
                func(cell *gtk.CellRendererText, path string, text string) {
                    iter, _ := show.treeList.GetIterFromString(path)
                    id, _ := show.treeList.GetValue(iter, PAGENUM)
                    val, _ := id.GoValue()
                    pageNum, _ := strconv.Atoi(val.(string))
                    show.pages[pageNum].title = text
                    show.treeList.SetValue(iter, TITLE, text)
            })

            column, _ = gtk.TreeViewColumnNewWithAttribute(KEYTITLE[key], title, "text", 1)
        default:
            cell, _ := gtk.CellRendererTextNew()
            column, _ = gtk.TreeViewColumnNewWithAttribute(KEYTITLE[key], cell, "text", key)
        }

        show.AppendColumn(column)
    }

    show.treeList, _ = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)

    show.SetModel(show.treeList)
    show.edit = edit

    // TODO: remove reference to show from outside scope
    show.Connect("row-activated", 
        func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) { 
            iter, _ := show.treeList.GetIter(path)
            id, _ := show.treeList.GetValue(iter, PAGENUM)
            val, _ := id.GoValue()
            pageNum, _ := strconv.Atoi(val.(string))

            show.edit.SetPage(show.pages[pageNum])
            conn["Preview"].SendPage(show.pages[pageNum], ANIMATE_ON)
        })

    return show 
}

func (show *ShowTree) NewShowPage(temp *Template) *Page {
    show.numPages++
    show.pages[show.numPages] = NewPage(show.numPages, temp.title, temp)
    page := show.pages[show.numPages]
    show.treeList.Set(
        show.treeList.Append(), 
        []int{PAGENUM, TITLE, TEMPLATEID}, 
        []interface{}{page.pageNum, page.title, page.templateID})
    return page
}

func (show *ShowTree) ImportShow(temp *TempTree, filename string) {
    pageReg, err := regexp.Compile("temp (?P<tempID>[0-9]*); title \"(?P<title>.*)\";")
    if err != nil {
        log.Print(err)
        return
    }

    propReg, err := regexp.Compile("name (?P<type>[ \\w]*);")
    if err != nil {
        log.Print(err)
        return
    }

    file, err := os.Open(filename)
    if err != nil {
        log.Print(err)
        return
    }

    scanner := bufio.NewScanner(file)

    var page *Page
    for scanner.Scan() {
        line := scanner.Text()
        if pageReg.Match(scanner.Bytes()) {
            match := pageReg.FindStringSubmatch(line)
            tempID, _ := strconv.Atoi(match[1])
            page = show.NewShowPage(temp.temps[tempID])
            page.title = match[2]
        } else if page != nil {
            match := propReg.FindStringSubmatch(line)
            if len(match) < 2 {
                log.Printf("Incorrect prop format (%s)\n", line)
                continue
            }

            name := match[1]

            if _, ok := page.props[name]; !ok {
                log.Printf("Unknown property (%s)\n", name)
                continue
            }

            page.props[name].Decode(line)
        }
    }
}

func (show *ShowTree) ExportShow(filename string) {
    file, err := os.Create(filename)
    if err != nil {
        log.Print(err)
    }
    defer file.Close()

    for _, page := range show.pages {
        pageString := fmt.Sprintf("temp %d; title \"%s\";\n", page.templateID, page.title)
        file.Write([]byte(pageString))

        for name, prop := range page.props {
            file.WriteString(fmt.Sprintf("name %s;", name))

            file.WriteString(prop.Encode())

            file.WriteString("\n")
        }
    }
}
