/* 
 * show_io.c 
 */ 

#include "chroma-viz.h"
#include "graphic.h"
#include <stdio.h>
#include <string.h>

#define TEXT_STRING           ""
#define RECTANGLE_STRING      "\t\t<rectangle; pos_x = %d; pos_y = %d; width = %d; height = %d; color = (%d, %d, %d, %d)>\n"
#define PAGE_STRING           "\t<page;prev = %s ;page_num = %d;title = %s ;num_rectangles = %d;num_text = %d>\n"
#define SHOW_STRING           "<show;show_name = %s ;num_pages = %d>\n"

enum {
    KEY_SHOW,
    KEY_PAGE,
    KEY_PAGE_ON,
    KEY_PAGE_OFF,
    KEY_RECTANGLE,
    KEY_TEXT
};

void write_page_to_file(FILE *, Graphic *, Page *);
void write_rectangle_to_file(FILE *, Rect *);
void write_text_to_file(FILE *, Text *);
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
        write_page_to_file(file, &show->graphic[i], &show->pages[i]);
    }

    fclose(file);
}

void write_page_to_file(FILE *file, Graphic *graphic, Page *page) {
    char buf[MAX_BUF_SIZE];

    memset(buf, '\0', sizeof buf);
    sprintf(buf, PAGE_STRING, page->prev, page->page_num, page->title, graphic->num_rect, graphic->num_text);
    fwrite(buf, 1, strlen(buf), file);

    for (int i = 0; i < graphic->num_rect; i++) {
        write_rectangle_to_file(file, &graphic->rect[i]);
    }

    for (int i = 0; i < graphic->num_text; i++) {
        write_text_to_file(file, &graphic->text[i]);
    }
}

void write_rectangle_to_file(FILE *file, Rect *rectangle) {
    char buf[MAX_BUF_SIZE];
    memset(buf, '\0', sizeof buf);

    sprintf(buf, RECTANGLE_STRING, 
            rectangle->pos_x, rectangle->pos_y, rectangle->width, rectangle->height, 
            rectangle->color.r, rectangle->color.g, rectangle->color.b, rectangle->color.a);

    fwrite(buf, 1, strlen(buf), file);
}

void write_text_to_file(FILE *file, Text *text) {

}

SHOW *read_show_from_file(char *filename) {
    FILE *file = fopen(filename, "r"); 
    SHOW *show = NEW_STRUCT( SHOW );
    Rect *rectangle;
    Text *text;

    char line[MAX_BUF_SIZE]; 
    int type, page_num = -1, rectangle_num, text_num;


    memset(line, '\0', sizeof line);

    while (fgets(line, MAX_BUF_SIZE, file)) {
        type = parse_line(line);

        if (type == KEY_SHOW) {
            show->show_name = NEW_ARRAY(MAX_CHAR_SIZE, char);
            memset(show->show_name, '\0', MAX_CHAR_SIZE);

            sscanf(line, SHOW_STRING, show->show_name, &show->num_pages);

            show->pages = NEW_ARRAY(show->num_pages, Page);
            show->graphic = NEW_ARRAY(show->num_pages, Graphic);
            show->page_start = 0;
            show->page_height = 30;

            // printf("%s, %d\n", show->show_name, show->num_pages);
        } else if (type == KEY_PAGE) {
            page_num++;
            rectangle_num = 0;
            text_num = 0;
            show->pages[page_num].prev  = NEW_ARRAY(MAX_CHAR_SIZE, char);
            show->pages[page_num].title = NEW_ARRAY(MAX_CHAR_SIZE, char);

            memset(show->pages[page_num].prev, '\0', MAX_CHAR_SIZE);
            memset(show->pages[page_num].title, '\0', MAX_CHAR_SIZE);

            sscanf(line, PAGE_STRING, 
                   show->pages[page_num].prev, &show->pages[page_num].page_num, show->pages[page_num].title, 
                   &show->graphic[page_num].num_rect, &show->graphic[page_num].num_text);

            //printf("%s, %d\n", show->pages[page_num].title, show->pages[page_num].page_num);
        } else if (type == KEY_RECTANGLE) {
            if (rectangle_num > sizeof(show->graphic[page_num]) % sizeof(Rect)) {
                printf("Too many rectangles\n");
            }

            rectangle = &show->graphic[page_num].rect[rectangle_num];

            sscanf(line, RECTANGLE_STRING, 
                   &rectangle->pos_x, &rectangle->pos_y, &rectangle->width, &rectangle->height, 
                   (int *) &rectangle->color.r , (int *) &rectangle->color.g , 
                   (int *) &rectangle->color.b, (int *) &rectangle->color.a);

            //printf("%d, %d, %d, %d\n", rectangle->pos_x, rectangle->pos_y, rectangle->width, rectangle->height);
            rectangle_num++;
        } else if (type == KEY_TEXT) {
            text_num++;
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
    } else if (strcmp(key, "rectangle") == 0) {
        return KEY_RECTANGLE;
    } else if (strcmp(key, "text") == 0) {
        return KEY_TEXT;
    }

    return -1;
}

