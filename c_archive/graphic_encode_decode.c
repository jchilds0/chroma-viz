/*
 * graphic_encode_decode.c 
 */ 

#include "chroma-prototypes.h"
#include "graphic.h"
#include <stdio.h>

#define NAME              "%s %c"
#define RECT_STRING       "pos_x=%d;pos_y=%d;width=%d;height=%d;color=(%d,%d,%d,%d)%c"
#define TEXT_STRING       "pos_x=%d;pos_y=%d;font_size=%d;text=%s;color=(%d,%d,%d,%d)%c"

void graphic_to_engine(int socket_desc, Graphic *page, bool animate_on) {
    char output[MAX_BUF_SIZE];
    Rect *rect;
    Text *text;

    memset(output, '\0', sizeof output);
    sprintf(output, NAME, page->name, END_OF_NAME);
    send_message_to_engine(socket_desc, output);

    for (int i = 0; i < page->num_rect; i++) {
        memset(output, '\0', sizeof output);
        rect = &page->rect[i];
        sprintf(output, RECT_STRING, rect->pos_x, rect->pos_y, rect->width, rect->height, 
                rect->color.r, rect->color.g, rect->color.b, rect->color.a, END_OF_RECTANGLE);
        send_message_to_engine(socket_desc, output);
    }
    
    for (int i = 0; i < page->num_text; i++) {
        memset(output, '\0', sizeof output);
        text = &page->text[i];
        sprintf(output, TEXT_STRING, text->pos_x, text->pos_y, text->font_size, text->text, 
                text->color.r, text->color.g, text->color.b, text->color.a, END_OF_TEXT);
        send_message_to_engine(socket_desc, output);
    }

    memset(output, '\0', sizeof output);
    output[0] = END_OF_GRAPHICS;
    send_message_to_engine(socket_desc, output);
}

void string_to_rectangle(Rect *, char *);
void string_to_text(Text *, char *);

void string_to_graphic(Graphic *graphic, char *buf) {
    int i = -1;
    int obj_index = -1;
    char obj[MAX_BUF_SIZE];
    memset(obj, '\0', sizeof obj);

    while ((obj[obj_index++] = buf[i++]) != END_OF_GRAPHICS) {
        if (buf[i] == END_OF_NAME) {
            strcpy(graphic->name, obj);

        } else if (buf[i] == END_OF_RECTANGLE) {
            string_to_rectangle(&graphic->rect[graphic->num_rect], obj);
            graphic->num_rect++;

        } else if (buf[i] == END_OF_TEXT) {
            string_to_text(&graphic->text[graphic->num_text], obj);
            graphic->num_text++;
        } else {
            continue;
        }

        memset(obj, '\0', sizeof obj);
        obj_index = 0;
    }
}

void string_to_rectangle(Rect *rect, char *buf) {
    char temp[MAX_BUF_SIZE];

    sscanf(buf, RECT_STRING, &rect->pos_x, &rect->pos_y, &rect->width, &rect->height, 
           (int *) &rect->color.r, (int *) &rect->color.g, (int *) &rect->color.b, (int *) &rect->color.a, temp);
}

void string_to_text(Text *text, char *buf) {
    char temp[MAX_BUF_SIZE];

    sscanf(buf, TEXT_STRING, &text->pos_x, &text->pos_y, &text->font_size, text->text, 
           (int *) &text->color.r, (int *) &text->color.g, (int *) &text->color.b, (int *) &text->color.a, temp);
}
