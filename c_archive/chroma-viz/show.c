/*
 * show.c 
 */

#include "chroma-typedefs.h"
#include "chroma-viz.h" 
#include <string.h>

void draw_header(TILE *);
void draw_page(TILE *, Graphic *, Page *, int, int, bool);
    
typedef struct {
    int       height;
    int       prev_width;
    int       title_width;
    int       page_num_width;
    int       x_pad_text;
    int       y_pad_text;
} HEADER;

static HEADER header = {25, 100, 400, 550, 10, 5};

SHOW *init_show(void) {
    SHOW *show = NEW_STRUCT( SHOW );

    show->show_name         = "basic_show";
    show->num_pages         = 2;
    show->page_start        = 0;
    show->selected_page     = 0;

    show->graphic           = NEW_ARRAY(show->num_pages, Graphic);
    show->pages             = NEW_ARRAY(show->num_pages, Page);
    show->page_height       = 30;

    show->pages[0]            = (Page) {"prev", "Green-Rectangle", 0};
    show->graphic[0].num_rect = 1; 
    show->graphic[0].num_text = 0; 
    show->graphic[0].rect[0]  = (Rect) {20, 20, 350, 100, GREEN};

    show->pages[1]            = (Page) {"prev", "Blue-Rectangle", 1};
    show->graphic[1].num_rect = 1; 
    show->graphic[1].num_text = 0; 
    show->graphic[1].rect[0]  = (Rect) {25, 500, 900, 200, BLUE};

    return show;
}

void free_show(SHOW *show) {
    free(show->pages);
    free(show);
}

void draw_show(TILE *show_tile, SHOW *show) {

    DrawRectangle(show_tile->pos_x, show_tile->pos_y, show_tile->width, show_tile->height, CHROMA_BG);
    //DrawText("Show", CENTER(show->pos_x, show->width), CENTER(show->pos_y, show->height), 20, CHROMA_TEXT);

    draw_header(show_tile);
    draw_page(show_tile, &show->graphic[0], &show->pages[0], header.height, show->page_height, 0 == show->selected_page);
    draw_page(show_tile, &show->graphic[1], &show->pages[1], header.height + show->page_height, show->page_height, 1 == show->selected_page);
}

void draw_header(TILE *show) {
    const int lower_pad = 3;
    const int div_width = 5;

    DrawRectangle(show->pos_x, show->pos_y, show->width, header.height, LIGHTGRAY);
    DrawRectangle(show->pos_x, show->pos_y + header.height - lower_pad, show->width, lower_pad, GRAY);
    DrawRectangle(show->pos_x, show->pos_y, div_width, header.height, GRAY);
    DrawRectangle(show->pos_x + header.prev_width, show->pos_y, div_width, header.height, GRAY);
    DrawRectangle(show->pos_x + header.title_width, show->pos_y, div_width, header.height, GRAY);
    DrawRectangle(show->pos_x + header.page_num_width, show->pos_y, div_width, header.height, GRAY);

    DrawText("Preview",     show->pos_x + header.x_pad_text,                       show->pos_y + header.y_pad_text, 20, BLACK);
    DrawText("Description", show->pos_x + header.prev_width + header.x_pad_text,  show->pos_y + header.y_pad_text, 20, BLACK);
    DrawText("Page Number", show->pos_x + header.title_width + header.x_pad_text, show->pos_y + header.y_pad_text, 20, BLACK);
}

void draw_page(TILE *show, Graphic *graphic, Page *page, int pos_y, int page_height, bool selected) {
    const int font_size = 20;
    const int x_pad = 10;
    const int y_pad = 10;
    char page_num[5];
    int prev_offset, title_offset, num_offset;

    memset(page_num, '\0', 5);
    sprintf(page_num, "%d", page->page_num);

    // draw bounding box
    DrawRectangle(show->pos_x, show->pos_y + pos_y, show->width, page_height, selected ? YELLOW : LIGHTGRAY);
    DrawRectangle(show->pos_x, show->pos_y + pos_y + page_height - 3, show->width, 3, GRAY);

    DrawRectangle(show->pos_x,                          show->pos_y + pos_y, 5, page_height, GRAY);
    DrawRectangle(show->pos_x + header.prev_width,     show->pos_y + pos_y, 5, page_height, GRAY);
    DrawRectangle(show->pos_x + header.title_width,    show->pos_y + pos_y, 5, page_height, GRAY);
    DrawRectangle(show->pos_x + header.page_num_width, show->pos_y + pos_y, 5, page_height, GRAY);

    // draw prev 
    if (page->prev != NULL) {

    }

    // draw title 
    title_offset = header.prev_width + x_pad; //header->title_width - MeasureText(page->title, font_size) - x_pad;
    DrawText(page->title, show->pos_x + title_offset,  show->pos_y + pos_y + y_pad, font_size, BLACK);

    // draw page_num
    num_offset = header.page_num_width - MeasureText(page_num, font_size) - x_pad;
    DrawText(page_num,    show->pos_x + num_offset, show->pos_y + pos_y + y_pad, font_size, BLACK);
}

void show_mouse_click(TILE *show_tile, SHOW *show, Connection *conn) {
    printf("Show mouse click\n");
    int page_index = -1;

    if (show_tile->pos_y + header.height <= GetMouseY())
        page_index = show->page_start + (GetMouseY() - header.height - show_tile->pos_y) / show->page_height;

    if (0 <= page_index && page_index < show->num_pages) {
        show->selected_page = page_index;
    }
}

