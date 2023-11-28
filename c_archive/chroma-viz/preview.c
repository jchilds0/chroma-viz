/*
 * preview.c 
 */

#include "chroma-viz.h" 
#include "graphic.h"
#include <raylib.h>

void draw_preview(TILE *preview, Connection *conn, Graphic *current_page) {
    Rect *rectangle;
    DrawRectangle(preview->pos_x, preview->pos_y, preview->width, preview->height, BLACK);
    //DrawText("Preview", CENTER(preview->pos_x, preview->width), CENTER(preview->pos_y, preview->height), 20, CHROMA_TEXT);

    for (int i = 0; i < current_page->num_rect; i++) {
        rectangle = &current_page->rect[i];
        DrawRectangle(preview->pos_x + rectangle->pos_x, preview->pos_y + rectangle->pos_y, 
                      rectangle->width * preview->width / conn->width, 
                      rectangle->height * preview->height / conn->height, 
                      rectangle->color);
    }
}

