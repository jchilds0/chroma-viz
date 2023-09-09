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

typedef struct {
    char *prev;
    int page_num;
    char *title;
} PAGE;

typedef struct {
    int height;
    int prev_width;
    int title_width;
    int page_num_width;
    int x_pad_text;
    int y_pad_text;
} HEADER;

void draw_header(TILE *, HEADER *);
void draw_page(TILE *, HEADER *, PAGE *, int, int);

void draw_show(TILE *show) {
    const int page_height = 30;
    HEADER header = (HEADER) {25, 100, 400, 550, 10, 5};
    PAGE page1 = (PAGE) {NULL, 001, "Green Rectangle"};

    DrawRectangle(show->pos_x, show->pos_y, show->width, show->height, CHROMA_BG);
    //DrawText("Show", CENTER(show->pos_x, show->width), CENTER(show->pos_y, show->height), 20, CHROMA_TEXT);

    draw_header(show, &header);
    draw_page(show, &header, &page1, header.height, page_height);
}

void draw_header(TILE *show, HEADER *header) {
    DrawRectangle(show->pos_x, show->pos_y, show->width, header->height, LIGHTGRAY);
    DrawRectangle(show->pos_x, show->pos_y + header->height - 3, show->width, 3, GRAY);
    DrawRectangle(show->pos_x, show->pos_y, 5, header->height, GRAY);
    DrawRectangle(show->pos_x + header->prev_width, show->pos_y, 5, header->height, GRAY);
    DrawRectangle(show->pos_x + header->title_width, show->pos_y, 5, header->height, GRAY);
    DrawRectangle(show->pos_x + header->page_num_width, show->pos_y, 5, header->height, GRAY);

    DrawText("Preview",     show->pos_x + header->x_pad_text,                      show->pos_y + header->y_pad_text, 20, BLACK);
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

void show_mouse_click(TILE *show, int socket_engine, int engine_status) {
    printf("Show mouse click\n");

    if (WITHIN(GetMouseY(), show->pos_y + 30, show->pos_y + 60) && engine_status == ENGINE_CON) {
        printf("Render Rectangle\n");
        RenderObject object = (RenderObject) {"rectangle", 50, 50, 10, 3, GREEN};
        render_objects(socket_engine, &object, 1);
    }
}
