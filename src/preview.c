/*
 * preview.c 
 */

#include "chroma-viz.h" 

void draw_preview(int pos_x, int pos_y, int width, int height) {
    DrawRectangle(pos_x, pos_y, width, height, CHROMA_BG);
    DrawText("Preview", CENTER(pos_x, width), CENTER(pos_y, height), 20, CHROMA_TEXT);
}
