/* 
 * show_io.c 
 */ 

#include "chroma-typedefs.h"
#include "chroma-viz.h"
#include "render_graphics.h"
#include <raylib.h>
#include <stdio.h>
#include <string.h>

#define MAX_BUF_SIZE          1000
#define MAX_CHAR_SIZE         50 

#define TEXT_STRING           ""
#define RECTANGLE_STRING      "\t\t\t<rectangle; pos_x = %d; pos_y = %d; width = %d; height = %d; color = (%d, %d, %d, %d)>\n"
#define PAGE_STRING           "\t<page;prev = %s ;page_num = %d;title = %s ;num_rectangles = %d ;num_text = %d>\n"
#define SHOW_STRING           "<show;show_name = %s ;num_pages = %d>\n"

enum {
    KEY_SHOW,
    KEY_PAGE,
    KEY_PAGE_ON,
    KEY_PAGE_OFF,
    KEY_RECTANGLE,
    KEY_TEXT
};

void write_page_to_file(FILE *, PAGE *);
void write_rectangle_to_file(FILE *, ChromaRectangle *);
void write_text_to_file(FILE *, ChromaText *);
int parse_line(char *);

void write_show_to_file(SHOW *show) {
    char filename[strlen(show->show_name) + 10]; 
    char buf[MAX_BUF_SIZE];

    memset(filename, '\0', sizeof filename);
    memset(buf, '\0', sizeof buf);

    sprintf(filename, "shows/%s.chromashow", show->show_name);
    remove(filename);
    FILE *file = fopen(filename, "w");

    sprintf(buf, SHOW_STRING, show->show_name, show->num_pages);
    fwrite(buf, 1, strlen(buf), file);

    for (int i = 0; i < show->num_pages; i++) {
        write_page_to_file(file, &show->pages[i]);
    }

    fclose(file);
}

void write_page_to_file(FILE *file, PAGE *page) {
    char buf[MAX_BUF_SIZE];

    memset(buf, '\0', sizeof buf);
    sprintf(buf, PAGE_STRING, page->prev, page->page_num, page->title, page->num_rectangles, page->num_text);
    fwrite(buf, 1, strlen(buf), file);

    memset(buf, '\0', sizeof buf);
    sprintf(buf, "\t\t<page_on>\n");
    fwrite(buf, 1, strlen(buf), file);

    for (int i = 0; i < page->num_rectangles; i++) {
        write_rectangle_to_file(file, &page->rectangles_on[i]);
    }

    for (int i = 0; i < page->num_text; i++) {
        write_text_to_file(file, &page->text_on[i]);
    }

    memset(buf, '\0', sizeof buf);
    sprintf(buf, "\t\t<page_off>\n");
    fwrite(buf, 1, strlen(buf), file);

    for (int i = 0; i < page->num_rectangles; i++) {
        write_rectangle_to_file(file, &page->rectangles_off[i]);
    }

    for (int i = 0; i < page->num_text; i++) {
        write_text_to_file(file, &page->text_off[i]);
    }
}

void write_rectangle_to_file(FILE *file, ChromaRectangle *rectangle) {
    char buf[MAX_BUF_SIZE];
    memset(buf, '\0', sizeof buf);

    sprintf(buf, RECTANGLE_STRING, 
            rectangle->pos_x, rectangle->pos_y, rectangle->width, rectangle->height, 
            rectangle->color.r, rectangle->color.g, rectangle->color.b, rectangle->color.a);

    fwrite(buf, 1, strlen(buf), file);
}

void write_text_to_file(FILE *file, ChromaText *text) {

}

SHOW *read_show_from_file(char *filename) {
    FILE *file = fopen(filename, "r"); 
    SHOW *show = NEW_STRUCT( SHOW );
    ChromaRectangle *rectangle;
    ChromaText *text;

    char line[MAX_BUF_SIZE]; 
    bool page_on;
    int type, page_num = -1, rectangle_num, text_num;


    memset(line, '\0', sizeof line);

    while (fgets(line, MAX_BUF_SIZE, file)) {
        type = parse_line(line);

        if (type == KEY_SHOW) {
            show->show_name = NEW_ARRAY(MAX_CHAR_SIZE, char);
            memset(show->show_name, '\0', MAX_CHAR_SIZE);

            sscanf(line, SHOW_STRING, show->show_name, &show->num_pages);

            show->pages = NEW_ARRAY(show->num_pages, PAGE);
            show->page_start = 0;
            show->page_height = 30;

            // printf("%s, %d\n", show->show_name, show->num_pages);
        } else if (type == KEY_PAGE) {
            page_num++;
            show->pages[page_num].prev  = NEW_ARRAY(MAX_CHAR_SIZE, char);
            show->pages[page_num].title = NEW_ARRAY(MAX_CHAR_SIZE, char);

            memset(show->pages[page_num].prev, '\0', MAX_CHAR_SIZE);
            memset(show->pages[page_num].title, '\0', MAX_CHAR_SIZE);

            sscanf(line, PAGE_STRING, show->pages[page_num].prev, &show->pages[page_num].page_num, show->pages[page_num].title, &show->pages[page_num].num_rectangles, &show->pages[page_num].num_text);

            show->pages[page_num].rectangles_on  = NEW_ARRAY(show->pages[page_num].num_rectangles, ChromaRectangle);
            show->pages[page_num].rectangles_off = NEW_ARRAY(show->pages[page_num].num_rectangles, ChromaRectangle);

            show->pages[page_num].text_on  = NEW_ARRAY(show->pages[page_num].num_text, ChromaText);
            show->pages[page_num].text_off = NEW_ARRAY(show->pages[page_num].num_text, ChromaText);

            //printf("%s, %d\n", show->pages[page_num].title, show->pages[page_num].page_num);
        } else if (type == KEY_PAGE_ON) {
            page_on = true;
            rectangle_num = 0;
            text_num = 0;

        } else if (type == KEY_PAGE_OFF) {
            page_on = false;
            rectangle_num = 0;
            text_num = 0;

        } else if (type == KEY_RECTANGLE) {
            if (page_on) {
                rectangle = &show->pages[page_num].rectangles_on[rectangle_num];
            } else {
                rectangle = &show->pages[page_num].rectangles_off[rectangle_num];
            }

            sscanf(line, RECTANGLE_STRING, 
                   &rectangle->pos_x, &rectangle->pos_y, &rectangle->width, &rectangle->height, 
                   &rectangle->color.r , &rectangle->color.g , &rectangle->color.b, &rectangle->color.a);

            //printf("%d, %d, %d, %d\n", rectangle->pos_x, rectangle->pos_y, rectangle->width, rectangle->height);
            rectangle_num++;
        } else if (type == KEY_TEXT) {
        }

        //printf("%d: %s", type, line);
    }

    fclose(file);
    return show;
}

int parse_line(char *buf) {
    char key[MAX_BUF_SIZE];
    int i = 0, j = 0;
    memset(key, '\0', sizeof key);

    while (buf[i++] != '<');
    while (buf[i] != ';' && buf[i] != '>') {
        key[j++] = buf[i++];
    }

    if (strcmp(key, "show") == 0) {
        return KEY_SHOW;
    } else if (strcmp(key, "page") == 0) {
        return KEY_PAGE;
    } else if (strcmp(key, "page_on") == 0) {
        return KEY_PAGE_ON;
    } else if (strcmp(key, "page_off") == 0) {
        return KEY_PAGE_OFF;
    } else if (strcmp(key, "rectangle") == 0) {
        return KEY_RECTANGLE;
    } else if (strcmp(key, "text") == 0) {
        return KEY_TEXT;
    }

    return -1;
}

