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

#endif // !CHROMA_TYPEDEFS
