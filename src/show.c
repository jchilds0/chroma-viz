/*
 * show.c 
 */

#include "chroma-prototypes.h"
#include "chroma-theme.h"
#include "chroma-typedefs.h"
#include "chroma-viz.h" 
#include <raylib.h>
#include <stdlib.h>
#include <string.h>


void draw_header(TILE *, HEADER *);
void draw_page(TILE *, HEADER *, PAGE *, int, int);

SHOW *init_show(void) {
    SHOW *show = NEW_STRUCT( SHOW );

    show->num_pages = 2;
    show->page_start = 0;
    show->pages = NEW_ARRAY(show->num_pages, PAGE);
    show->page_graphic = NEW_ARRAY(show->num_pages, RenderObject);
    show->page_height = 30;

    show->header = (HEADER) {25, 100, 400, 550, 10, 5};

    show->pages[0]        = (PAGE) {NULL, 001, "Green Rectangle On"};
    show->page_graphic[0] = (RenderObject) {"rectangle", 50, 50, 100, 30, GREEN};

    show->pages[1]        = (PAGE) {NULL, 001, "Green Rectangle Off"};
    show->page_graphic[1] = (RenderObject) {"rectangle", 50, 50, 100, 30, BLACK};

    return show;
}

void free_show(SHOW *show) {
    free(show->pages);
    free(show);
}

void draw_show(TILE *show_tile, SHOW *show) {

    DrawRectangle(show_tile->pos_x, show_tile->pos_y, show_tile->width, show_tile->height, CHROMA_BG);
    //DrawText("Show", CENTER(show->pos_x, show->width), CENTER(show->pos_y, show->height), 20, CHROMA_TEXT);

    draw_header(show_tile, &show->header);
    draw_page(show_tile, &show->header, &show->pages[0], show->header.height, show->page_height);
    draw_page(show_tile, &show->header, &show->pages[1], show->header.height + show->page_height, show->page_height);
}

void draw_header(TILE *show, HEADER *header) {
    const int lower_pad = 3;
    const int div_width = 5;

    DrawRectangle(show->pos_x, show->pos_y, show->width, header->height, LIGHTGRAY);
    DrawRectangle(show->pos_x, show->pos_y + header->height - lower_pad, show->width, lower_pad, GRAY);
    DrawRectangle(show->pos_x, show->pos_y, div_width, header->height, GRAY);
    DrawRectangle(show->pos_x + header->prev_width, show->pos_y, div_width, header->height, GRAY);
    DrawRectangle(show->pos_x + header->title_width, show->pos_y, div_width, header->height, GRAY);
    DrawRectangle(show->pos_x + header->page_num_width, show->pos_y, div_width, header->height, GRAY);

    DrawText("Preview",     show->pos_x + header->x_pad_text,                       show->pos_y + header->y_pad_text, 20, BLACK);
    DrawText("Description", show->pos_x + header->prev_width + header->x_pad_text,  show->pos_y + header->y_pad_text, 20, BLACK);
    DrawText("Page Number", show->pos_x + header->title_width + header->x_pad_text, show->pos_y + header->y_pad_text, 20, BLACK);
}

void draw_page(TILE *show, HEADER *header, PAGE *page, int pos_y, int page_height) {
    const int font_size = 20;
    const int x_pad = 10;
    const int y_pad = 10;
    char page_num[5];
    int prev_offset, title_offset, num_offset;

    memset(page_num, '\0', 5);
    sprintf(page_num, "%d", page->page_num);

    // draw bounding box
    DrawRectangle(show->pos_x, show->pos_y + pos_y, show->width, page_height, LIGHTGRAY);
    DrawRectangle(show->pos_x, show->pos_y + pos_y + page_height - 3, show->width, 3, GRAY);

    DrawRectangle(show->pos_x,                          show->pos_y + pos_y, 5, page_height, GRAY);
    DrawRectangle(show->pos_x + header->prev_width,     show->pos_y + pos_y, 5, page_height, GRAY);
    DrawRectangle(show->pos_x + header->title_width,    show->pos_y + pos_y, 5, page_height, GRAY);
    DrawRectangle(show->pos_x + header->page_num_width, show->pos_y + pos_y, 5, page_height, GRAY);

    // draw prev 
    if (page->prev != NULL) {

    }

    // draw title 
    title_offset = header->prev_width + x_pad; //header->title_width - MeasureText(page->title, font_size) - x_pad;
    DrawText(page->title, show->pos_x + title_offset,  show->pos_y + pos_y + y_pad, font_size, BLACK);

    // draw page_num
    num_offset = header->page_num_width - MeasureText(page_num, font_size) - x_pad;
    DrawText(page_num,    show->pos_x + num_offset, show->pos_y + pos_y + y_pad, font_size, BLACK);
}

void show_mouse_click(TILE *show_tile, SHOW *show, int socket_engine, int engine_status) {
    printf("Show mouse click\n");
    int page_index = -1;

    if (show_tile->pos_y + show->header.height <= GetMouseY())
        page_index = show->page_start + (GetMouseY() - show->header.height - show_tile->pos_y) / show->page_height;

    if (0 <= page_index && page_index < show->num_pages) {
        printf("Render Rectangle\n");

        if (engine_status == ENGINE_CON) {
            render_objects(socket_engine, &show->page_graphic[page_index], 1);
        }
    }
}
