/*
 * chroma-hub.h 
 */

#ifndef CHROMA_HUB
#define CHROMA_HUB

#include <raylib.h>

typedef struct {
    int           pos_x;
    int           pos_y;
    int           width;
    int           height;
} Rectangle;

typedef struct {
    int           font_size;
    char          text;
} Text;

typedef struct {
    int           num_rect;
    int           num_text;
    Rectangle     *rect;
    Text          *text;
} Graphic;

typedef struct {
    int           len;
    int           max_size;
    Graphic       *pages;
} HashTable;

#endif // !CHROMA_HUB
