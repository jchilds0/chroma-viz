/*
 * editor.c 
 */

#include "chroma-viz.h" 

void draw_editor(TILE *editor) {
    DrawRectangle(editor->pos_x, editor->pos_y, editor->width, editor->height, CHROMA_BG);
    DrawText("Editor", CENTER(editor->pos_x, editor->width), CENTER(editor->pos_y, editor->height), 20, CHROMA_TEXT);
}
