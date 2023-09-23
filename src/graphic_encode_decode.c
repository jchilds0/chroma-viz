/*
 * graphic_encode_decode.c 
 */ 

#include "graphic.h"

void graphic_to_string(Graphic *page, char *buf) {
    buf[0] = END_OF_GRAPHICS;
}

void string_to_rectangle(Rect *, char *);
void string_to_text(Text *, char *);

void string_to_graphic(Graphic *graphic, char *buf) {
    int i = -1;
    int obj_index = 0;
    char obj[MAX_BUF_SIZE];
    memset(obj, '\0', sizeof obj);

    while (buf[i++] != END_OF_GRAPHICS) {
        if (buf[i] == END_OF_NAME) {
            strcpy(graphic->name, obj);

        } else if (buf[i] == END_OF_RECTANGLE) {
            string_to_rectangle(&graphic->rect[graphic->num_rect], obj);
            graphic->num_rect++;

        } else if (buf[i] == END_OF_TEXT) {
            string_to_text(&graphic->text[graphic->num_text], obj);
            graphic->num_text++;

        } else {
            obj[obj_index++] = buf[i];
            continue;
        }

        memset(obj, '\0', sizeof obj);
        obj_index = 0;
    }
}

void string_to_rectangle(Rect *rect, char *buf) {
    sscanf(buf, "pos_x=%d;pos_y=%d;width=%d;height=%d;color=(%c,%c,%c,%c)", 
           &rect->pos_x, &rect->pos_y, &rect->width, &rect->height, &rect->color.r, &rect->color.g, &rect->color.b, &rect->color.a);
}

void string_to_text(Text *text, char *buf) {
    sscanf(buf, "pos_x=%d;pos_y=%d;font_size=%d;text=%s;color=(%c,%c,%c,%c)", 
           &text->pos_x, &text->pos_y, &text->font_size, text->text, 
           &text->color.r, &text->color.g, &text->color.b, &text->color.a);
}
