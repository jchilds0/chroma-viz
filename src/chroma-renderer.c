/*
 * chroma-renderer.c 
 */ 

#include "chroma-prototypes.h"
#include "chroma-typedefs.h"
#include "chroma-viz.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MAX_POS_LEN               4
#define MAX_COLOR_LEN             3
#define MAX_OUTPUT_WIDTH          1920
#define MAX_OUTPUT_HEIGHT         1080
#define PIXEL_LEN                 (2 * MAX_POS_LEN + 4 * MAX_COLOR_LEN + 10)
#define MAX_BUF_SIZE              100

typedef struct {
    int       pos_x;
    int       pos_y;
    Color     color;
} RenderPixel;

RenderPixel *rectangle_to_pixels(RenderObject *, int *);
void render_object_to_str(RenderPixel *, char *);

int render_objects(int socket_desc, RenderObject *objects, int num_objects) {
    int num_pixels;
    RenderObject *object;
    RenderPixel *pixels;
    char buf[MAX_BUF_SIZE];

    for (int i = 0; i < num_objects; i++) {
        object = &objects[i];

        if (strcmp(object->type, "rectangle") == 0) {
            pixels = rectangle_to_pixels(object, &num_pixels);
        }

        for (int j = 0; j < num_pixels; j++) {
            memset(buf, '\0', sizeof buf );
            render_object_to_str(&pixels[j], buf);
            send_message_to_engine(socket_desc, buf);
        }

        free(pixels);
    }

    memset(buf, '\0', sizeof buf );
    buf[0] = END_OF_FRAME;
    send_message_to_engine(socket_desc, buf);
    return 1;
}

RenderPixel *rectangle_to_pixels(RenderObject *object, int *num_pixels) {
    (*num_pixels) = object->width * object->height;
    RenderPixel *pixels = NEW_ARRAY(*num_pixels, RenderPixel);
    int i = 0;

    for (int x = object->pos_x; x < object->pos_x + object->width; x++) {
        for (int y = object->pos_y; y < object->pos_y + object->height; y++) {
            pixels[i].pos_x = x;
            pixels[i].pos_y = y;
            pixels[i].color = object->color;
            i++;
        }
    }

    return pixels;
}

void render_object_to_str(RenderPixel *pixel, char *buf) {
    sprintf(buf, "%d,%d,%d,%d,%d%c", pixel->pos_x, pixel->pos_y, 
            pixel->color.r, pixel->color.g, pixel->color.b, END_OF_PIXEL);

    //printf("%s\n", buf);
}
