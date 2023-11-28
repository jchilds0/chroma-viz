package gui

var show map[int]Page

NewShowPage(1, )

func NewShowPage(pageNum int, template Page) bool {
    if _, ok := show[pageNum]; !ok {
        return false
    }

    show[pageNum] = *NewPage(pageNum, template.title, template.templateID)
    return true
}



