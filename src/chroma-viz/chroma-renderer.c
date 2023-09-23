/*
 * chroma-renderer.c 
 */ 

#include "chroma-prototypes.h"
#include "chroma-viz.h"
#include "graphic.h"
#include <string.h>

#define MAX_POS_LEN               4
#define MAX_COLOR_LEN             3
#define MAX_OUTPUT_WIDTH          1920
#define MAX_OUTPUT_HEIGHT         1080
#define PIXEL_LEN                 (2 * MAX_POS_LEN + 4 * MAX_COLOR_LEN + 10)
#define MAX_BUF_SIZE              100

void render_graphic(Graphic *page, Connection *conn, bool animate_on) {
    char buf[MAX_BUF_SIZE];

    memset(buf, '\0', sizeof buf);
    graphic_to_string(page, buf);
    send_message_to_engine(conn->socket_desc, buf);
}
