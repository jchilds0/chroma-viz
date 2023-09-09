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

#define RENDERER_PIXELS               2
#define RENDERER_TEXT                 3

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

#endif // !CHROMA_TYPEDEFS
