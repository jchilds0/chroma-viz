/*
 * chroma-typedefs.h
 */ 

#ifndef CHROMA_TYPEDEFS
#define CHROMA_TYPEDEFS

#include "raylib.h"

#define CENTER(pos, offset)     (pos + offset / 2)

typedef struct {
    int pos_x;
    int pos_y;
    int width;
    int height;
    int split;
} PANE;

typedef struct {
    char      type[30];
    int       pos_x;
    int       pos_y;
    int       width;
    int       height;
} render_object;

#endif // !CHROMA_TYPEDEFS
