/*
 * chroma-typedefs.h
 */ 

#ifndef CHROMA_TYPEDEFS
#define CHROMA_TYPEDEFS

#include "raylib.h"
#include <malloc.h>

#define NEW_STRUCT(struct_type)       (struct_type *) malloc((size_t) sizeof( struct_type ) )
#define NEW_ARRAY(n, struct_type)     (struct_type *) malloc((size_t) (n) * sizeof( struct_type ))

#define WITHIN(x, x0, x1)             (x0 <= x && x <= x1)

#define CENTER(pos, offset)           (pos + offset / 2)
#define RENDER_NUM_PARAMS             4

#define RENDER_PIXELS                 2
#define RENDER_TEXT                   3
#define END_OF_PIXEL                  4
#define END_OF_FRAME                  5
#define END_OF_CON                    6

#define ENGINE_DISCON                 0
#define ENGINE_TEST                   1
#define ENGINE_CON                    2

typedef struct {
    int pos_x;
    int pos_y;
    int width;
    int height;
    int split;
} PANE;

typedef struct {
    int pos_x;
    int pos_y;
    int width;
    int height;
} TILE;

typedef struct {
    char      type[30];
    int       pos_x;
    int       pos_y;
    int       width;
    int       height;
    Color     color;
} RenderObject;

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

typedef struct {
    int       page_start;
    int       num_pages;
    int       page_height;
    PAGE      *pages;
    RenderObject *page_graphic;
    HEADER    header;
} SHOW;

#endif // !CHROMA_TYPEDEFS
