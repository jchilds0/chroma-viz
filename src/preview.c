/*
 * preview.c 
 */

#include "chroma-viz.h" 

void draw_preview(TILE *preview) {
    DrawRectangle(preview->pos_x, preview->pos_y, preview->width, preview->height, CHROMA_BG);
    DrawText("Preview", CENTER(preview->pos_x, preview->width), 
             CENTER(preview->pos_y, preview->height), 20, CHROMA_TEXT);
}
