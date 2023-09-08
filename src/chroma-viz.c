/*
 * chroma-viz.c 
 */

#include "chroma-viz.h" 
#include "chroma-prototypes.h"
#include <raylib.h>
#include <stdio.h>
#include <sys/socket.h>

void left_pane(PANE *, int);
void right_pane(PANE *, int);

int main(int argc, char **argv) {
    const int screen_width = 1600;
    const int screen_height = 1000;
    const int padding = 10;
    PANE main  = {0, 0, screen_width, screen_height, screen_width / 2};
    PANE left, right;
    int split_left = main.height / 2, split_right = main.height / 2 - 100;

    // Navbar
    main.pos_y = 20;
    main.height = main.height - main.pos_y;

    InitWindow(screen_width, screen_height, "raylib [core] example - basic window");
    SetTargetFPS(30);

    int socket_engine = connect_to_engine("127.0.0.1", 6100);
    int i = 0;

    while (!WindowShouldClose()) {
        BeginDrawing();
        ClearBackground(RAYWHITE);

        // Navbar
        DrawRectangle(0, 0, main.width, main.pos_y, YELLOW);

        char message[50];
        sprintf(message, "Time: %d", i++ / 10);
        if (send_message_to_engine(socket_engine, message) < 0) {
            printf("Error sending message to server\n");
        }

        left  = (PANE) {
            main.pos_x + padding, 
            main.pos_y + padding, 
            main.split - padding - padding / 2, 
            main.height - 2 * padding, 
            split_left
        };

        right = (PANE) {
            main.pos_x + main.split + padding / 2, 
            main.pos_y + padding, 
            main.width - main.split - padding - padding / 2, 
            main.height - 2 * padding, 
            split_right
        };

        left_pane(&left, padding);
        right_pane(&right, padding);

        EndDrawing();
    }

    close_engine_connection(socket_engine);

    CloseWindow();
    return 0;
}

void left_pane(PANE *pane, int padding) {
    draw_templates(pane->pos_x, pane->pos_y, pane->width, pane->split - padding / 2);
    draw_show(pane->pos_x, pane->pos_y + pane->split + padding / 2, pane->width, pane->height - pane->split - padding / 2);
}

void right_pane(PANE *pane, int padding) {
    draw_editor(pane->pos_x, pane->pos_y, pane->width, pane->split - padding / 2);
    draw_preview(pane->pos_x, pane->pos_y + pane->split + padding / 2, pane->width, pane->height - pane->split - padding / 2);
}
