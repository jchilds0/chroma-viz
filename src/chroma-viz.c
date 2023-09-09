/*
 * chroma-viz.c 
 */

#include "chroma-viz.h" 
#include "chroma-prototypes.h"
#include "chroma-typedefs.h"
#include <raylib.h>
#include <stdio.h>
#include <sys/socket.h>

#define PADDING     5

void left_pane(PANE *, TILE *, TILE *);
void right_pane(PANE *, TILE *, TILE *);
void update_window_panes(PANE *, PANE *, PANE *);
void update_pane_tiles(PANE *, TILE *, TILE *);
void draw_connection_button(PANE *, TILE *, bool);
void connection_mouse_click(int *, int *);

int main(int argc, char **argv) {
    const int screen_width = 1600;
    const int screen_height = 1000;
    PANE main  = {0, 0, screen_width, screen_height, screen_width / 2};
    PANE left, right;
    TILE editor, preview, templates, show;
    int engine_status = ENGINE_DISCON;
    int socket_engine = 0;

    // Navbar
    main.pos_y = 20;
    main.height = main.height - main.pos_y - 30;
    TILE lower = (TILE) {main.pos_x, main.pos_y + main.height, main.width, 30};

    InitWindow(screen_width, screen_height, "raylib [core] example - basic window");
    SetTargetFPS(30);

    while (!WindowShouldClose()) {
        BeginDrawing();
        ClearBackground(RAYWHITE);

        // Navbar
        DrawRectangle(0, 0, main.width, main.pos_y, YELLOW);

        // Lower bar
        if (IsMouseButtonPressed(MOUSE_BUTTON_LEFT)) { 
            if (WITHIN(GetMouseX(), lower.pos_x, lower.pos_x + 100)
                && WITHIN(GetMouseY(), lower.pos_y, lower.pos_y + 30)) {

                    connection_mouse_click(&socket_engine, &engine_status);

            } else if (WITHIN(GetMouseX(), show.pos_x, show.pos_x + show.width) 
                && WITHIN(GetMouseY(), show.pos_y, show.pos_y + show.height)) {

                    show_mouse_click(&show, socket_engine, engine_status);

            }
        }

        draw_connection_button(&main, &lower, engine_status);

        update_window_panes(&main, &left, &right);
        update_pane_tiles(&left, &templates, &show);
        update_pane_tiles(&right, &editor, &preview);

        left_pane(&left, &templates, &show);
        right_pane(&right, &editor, &preview);

        EndDrawing();
    }

    if (engine_status == ENGINE_CON) {
        close_engine_connection(socket_engine);
    }

    CloseWindow();
    return 0;
}

void update_window_panes(PANE *main, PANE *left, PANE *right) {
    left->pos_x = main->pos_x + PADDING; 
    left->pos_y = main->pos_y + PADDING; 
    left->width = main->split - PADDING - PADDING/ 2;
    left->height = main->height - 2 * PADDING;
    left->split = main->height / 2;

    right->pos_x = main->pos_x + main->split + PADDING / 2;
    right->pos_y = main->pos_y + PADDING;
    right->width = main->width - main->split - PADDING - PADDING / 2;
    right->height = main->height - 2 * PADDING;
    right->split = right->height - (right->split) * 9 / 16;
}

void update_pane_tiles(PANE *pane, TILE *top, TILE *bottom) {
    top->pos_x = pane->pos_x;
    top->pos_y = pane->pos_y;
    top->width = pane->width;
    top->height = pane->split - PADDING / 2;

    bottom->pos_x = pane->pos_x;
    bottom->pos_y = pane->pos_y + pane->split + PADDING / 2;
    bottom->width = pane->width;
    bottom->height = pane->height - pane->split - PADDING / 2;
}

void left_pane(PANE *left, TILE *templates, TILE *show) {
    draw_templates(templates);
    draw_show(show);
}

void right_pane(PANE *pane, TILE *editor, TILE *preview) {
    draw_editor(editor);
    draw_preview(preview);
}

void draw_connection_button(PANE *main, TILE *lower, bool status) {
    *lower = (TILE) {main->pos_x, main->pos_y + main->height, main->width, 30};

    DrawRectangle(lower->pos_x, lower->pos_y, lower->width, lower->height, WHITE);
    DrawRectangle(lower->pos_x, lower->pos_y, 100, lower->height, status ? GREEN : RED);

}

void connection_mouse_click(int *socket_engine, int *engine_status) {
    if (*engine_status == ENGINE_DISCON) {
    // draw_connection_button(&main, &lower, ENGINE_TEST);
        *socket_engine = connect_to_engine("127.0.0.1", 6100);
        *engine_status = *socket_engine < 0 ? ENGINE_DISCON : ENGINE_CON;
    } else if (*engine_status == ENGINE_CON) {
        close_engine_connection(*socket_engine);
        *engine_status = ENGINE_DISCON;
    }
}
