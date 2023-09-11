/*
 * chroma-viz.c 
 */

#include "chroma-viz.h" 
#include "chroma-prototypes.h"
#include "chroma-typedefs.h"
#include <raylib.h>
#include <stdatomic.h>
#include <stdio.h>
#include <sys/socket.h>

#define PADDING     5

void left_pane(PANE *, TILE *, TILE *, SHOW *);
void right_pane(PANE *, TILE *, TILE *);
void update_window_panes(PANE *, PANE *, PANE *);
void update_pane_tiles(PANE *, TILE *, TILE *);
void update_lower_bar(PANE *, TILE *);
void draw_connection_button(PANE *, TILE *, bool);
void connection_mouse_click(int *, int *);

int main(int argc, char **argv) {
    PANE main  = {0, 0, 1200, 800, 600};
    PANE left, right;
    TILE editor, preview, templates, show_tile;
    int engine_status = ENGINE_DISCON;
    int socket_engine = 0;
    bool resize_mid = false, resize_left = false;

    // Navbar
    main.pos_y = 20;
    main.height = main.height - main.pos_y - 30;
    left.split = main.height / 2;
    TILE lower = (TILE) {main.pos_x, main.pos_y + main.height, main.width, 30};
    SHOW *show = init_show();

    SetWindowState(FLAG_WINDOW_RESIZABLE);
    InitWindow(main.width, main.height, "raylib [core] example - basic window");
    SetTargetFPS(30);

    while (!WindowShouldClose()) {
        BeginDrawing();
        ClearBackground(RAYWHITE);

        if (main.width != GetScreenWidth()) {
            main.split = GetScreenWidth() / 2;
        }
        main.width = GetScreenWidth();
        main.height = GetScreenHeight() - main.pos_y - lower.height;

        // Navbar
        DrawRectangle(0, 0, main.width, main.pos_y, YELLOW);

        // Lower bar
        if (IsMouseButtonPressed(MOUSE_BUTTON_LEFT)) { 
            if (WITHIN(GetMouseX(), lower.pos_x, lower.pos_x + 100)
                && WITHIN(GetMouseY(), lower.pos_y, lower.pos_y + 30)) {

                    connection_mouse_click(&socket_engine, &engine_status);

            } else if (WITHIN(GetMouseX(), show_tile.pos_x, show_tile.pos_x + show_tile.width) 
                && WITHIN(GetMouseY(), show_tile.pos_y, show_tile.pos_y + show_tile.height)) {

                    show_mouse_click(&show_tile, show, socket_engine, engine_status);
            }
        }
        
        if (IsMouseButtonDown(MOUSE_BUTTON_LEFT)) {
            if (WITHIN(GetMouseX(), main.pos_x + main.split - PADDING, main.pos_x + main.split + PADDING)) {
                resize_mid = true;
            } else if (WITHIN(GetMouseY(), left.pos_y + left.split - PADDING, left.pos_y + left.split + PADDING) 
                    && WITHIN(GetMouseX(), left.pos_x, main.pos_x + main.split)) {
                resize_left = true;
            }
        }

        if (IsMouseButtonUp(MOUSE_BUTTON_LEFT)) {
            resize_mid = false;
            resize_left = false;
        }

        if (resize_mid) {
            main.split = GetMouseX();
        } else if (resize_left) {
            left.split = GetMouseY() - left.pos_y;
        }

        draw_connection_button(&main, &lower, engine_status);

        update_window_panes(&main, &left, &right);
        update_pane_tiles(&left, &templates, &show_tile);
        update_pane_tiles(&right, &editor, &preview);
        update_lower_bar(&main, &lower);

        left_pane(&left, &templates, &show_tile, show);
        right_pane(&right, &editor, &preview);

        EndDrawing();
    }

    if (engine_status == ENGINE_CON) {
        close_engine_connection(socket_engine);
    }

    free_show(show);

    CloseWindow();
    return 0;
}

void update_window_panes(PANE *main, PANE *left, PANE *right) {
    left->pos_x = main->pos_x + PADDING; 
    left->pos_y = main->pos_y + PADDING; 
    left->width = main->split - PADDING - PADDING/ 2;
    left->height = main->height - 2 * PADDING;

    right->pos_x = main->pos_x + main->split + PADDING / 2;
    right->pos_y = main->pos_y + PADDING;
    right->width = main->width - main->split - PADDING - PADDING / 2;
    right->height = main->height - 2 * PADDING;
    right->split = right->height - (right->width) * 9 / 16;
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

void update_lower_bar(PANE *main, TILE *lower) {
    *lower = (TILE) {main->pos_x, main->pos_y + main->height, main->width, 30};
}

void left_pane(PANE *left, TILE *templates, TILE *show_tile, SHOW *show) {
    draw_templates(templates);
    draw_show(show_tile, show);
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
