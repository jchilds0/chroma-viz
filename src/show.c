/*
 * show.c 
 */

#include "chroma-viz.h" 

void draw_show(int pos_x, int pos_y, int width, int height) {
    DrawRectangle(pos_x, pos_y, width, height, CHROMA_BG);
    DrawText("Show", CENTER(pos_x, width), CENTER(pos_y, height), 20, CHROMA_TEXT);
}
