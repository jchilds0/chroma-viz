/*
 * lower_bar.c 
 */

#include "chroma-typedefs.h"
#include "chroma-viz.h"
#include <raylib.h>

void draw_connection_button(PANE *main, TILE *lower, Connection *engine, Connection *prev) {
    *lower = (TILE) {main->pos_x, main->pos_y + main->height, main->width, 30};

    DrawRectangle(lower->pos_x, lower->pos_y, lower->width, lower->height, WHITE);
    DrawRectangle(lower->pos_x, lower->pos_y, 100, lower->height, engine->status ? GREEN : RED);
    DrawRectangle(lower->pos_x + 110, lower->pos_y, 100, lower->height, prev->status ? GREEN : RED);
}

void update_lower_bar(PANE *main, TILE *lower) {
    *lower = (TILE) {main->pos_x, main->pos_y + main->height, main->width, 30};
}

void lower_bar_mouse_click(TILE *lower, Connection *engine, Connection *prev) {
    if (WITHIN(GetMouseX(), lower->pos_x, lower->pos_x + 100)) {
        if (engine->status == ENGINE_DISCON) {
            engine->socket_desc = connect_to_engine(engine->addr, engine->port);
            engine->status = engine->socket_desc < 0 ? ENGINE_DISCON : ENGINE_CON;

        } else if (engine->status == ENGINE_CON) {
            close_engine_connection(engine->socket_desc);
            engine->status = ENGINE_DISCON;

        }

    } else if (WITHIN(GetMouseX(), lower->pos_x + 110, lower->pos_y + 210)) {
        if (prev->status == ENGINE_DISCON) {
            prev->socket_desc = connect_to_engine(prev->addr, prev->port);
            prev->status = prev->socket_desc < 0 ? ENGINE_DISCON : ENGINE_CON;

        } else if (prev->status == ENGINE_CON) {
            close_engine_connection(engine->socket_desc);
            engine->status = ENGINE_DISCON;

        }

    }
}
