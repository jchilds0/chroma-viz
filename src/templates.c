/*
 * templates.c 
 */

#include "chroma-theme.h"
#include "chroma-viz.h" 

void draw_templates(TILE *templates) {
    DrawRectangle(templates->pos_x, templates->pos_y, templates->width, templates->height, CHROMA_BG);
    DrawText("Templates", CENTER(templates->pos_x, templates->width), 
             CENTER(templates->pos_y, templates->height), 20, CHROMA_TEXT);
}

