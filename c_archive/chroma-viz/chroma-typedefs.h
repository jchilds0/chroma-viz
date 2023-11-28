/*
 * chroma-typedefs.h
 */ 

#ifndef CHROMA_TYPEDEFS
#define CHROMA_TYPEDEFS

#include "raylib.h"
#include <malloc.h>
#include "graphic.h"

#define NEW_STRUCT(struct_type)       (struct_type *) malloc((size_t) sizeof( struct_type ) )
#define NEW_ARRAY(n, struct_type)     (struct_type *) malloc((size_t) (n) * sizeof( struct_type ))

#define WITHIN(x, x0, x1)             (x0 <= x && x <= x1)

#define CENTER(pos, offset)           (pos + offset / 2)
#define RENDER_NUM_PARAMS             4

/*
 * main display components
 */

typedef struct {
    int       pos_x;
    int       pos_y;
    int       width;
    int       height;
    int       split;
} PANE;

typedef struct {
    int       pos_x;
    int       pos_y;
    int       width;
    int       height;
} TILE;

typedef struct {
    int       width;
    int       height;
    int       status;
    int       socket_desc;
    char      *addr;
    int       port;
} Connection;

/*
 * show templates
 */

typedef struct {
    char            *prev;
    char            *title;
    int             page_num;
} Page;

typedef struct {
    int             page_start;
    int             num_pages;
    int             page_height;
    int             selected_page;
    char            *show_name;
    Graphic         *graphic;
    Page            *pages;
} SHOW;

/*
 * keymap object 
 */ 

typedef struct {
    int       animate_on;
    int       animate_off;
    int       next_page;
    int       prev_page;
} Keymap;

#endif // !CHROMA_TYPEDEFS
