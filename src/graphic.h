/* 
 * graphic.h 
 */

#ifndef CHROMA_GRAPHICS
#define CHROMA_GRAPHICS

#include <raylib.h>
#include <stdio.h>
#include <string.h>

#define MAX_BUF_SIZE                  100
#define MAX_CHAR_SIZE                 100
#define MAX_GRAPHICS                  30
#define MAX_SUB_COMPS                 10

#define END_OF_CON                    2
#define END_OF_GRAPHICS               3
#define END_OF_NAME                   4
#define END_OF_RECTANGLE              5
#define END_OF_TEXT                   6

typedef struct {
    int             pos_x;
    int             pos_y;
    int             width;
    int             height;
    Color           color;
} Rect;

typedef struct {
    int             pos_x;
    int             pos_y;
    int             font_size;
    char            text[MAX_CHAR_SIZE];
    Color           color;
} Text;

typedef struct {
    char            name[MAX_CHAR_SIZE];
    int             num_rect;
    int             num_text;
    Rect            rect[MAX_SUB_COMPS];
    Text            text[MAX_SUB_COMPS];
} Graphic;

void graphic_to_string(Graphic *, char *);
void string_to_graphic(Graphic *, char *);

#endif // !CHROMA_GRAPHICS

