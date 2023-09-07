/*
 * chroma-viz.c 
 */

#include "chroma-viz.h" 

void left_pane(PANE *, int);
void right_pane(PANE *, int);

int main(int argc, char **argv) {
    const int screen_width = 1600;
    const int screen_height = 1000;
    const int padding = 10;
    PANE main  = {0, 0, screen_width, screen_height, screen_width / 2};
    PANE left, right;
    int split_left = main.height / 2, split_right = main.height / 2 - 100;

    InitWindow(screen_width, screen_height, "raylib [core] example - basic window");
    SetTargetFPS(60);

    while (!WindowShouldClose()) {
        BeginDrawing();

        ClearBackground(RAYWHITE);

        left  = (PANE) {
            main.pos_x, 
            main.pos_y, 
            main.split, 
            main.height, 
            split_left
        };

        right = (PANE) {
            main.pos_x + main.split, 
            main.pos_y, 
            main.width - main.split, 
            main.height, 
            split_right
        };

        left_pane(&left, padding);
        right_pane(&right, padding);

        EndDrawing();
    }

    CloseWindow();
    return 0;
}

void left_pane(PANE *pane, int padding) {
    int width = pane->width - 3 * padding / 2;
    int height = pane->height - 2 * padding;

    draw_templates(pane->pos_x + padding, pane->pos_y + padding, width, pane->split - padding);
    draw_show(pane->pos_x + padding, pane->pos_y + pane->split + padding, width, height - pane->split - padding);
}

void right_pane(PANE *pane, int padding) {
    int width = pane->width - 3 * padding / 2;
    int height = pane->height - 2 * padding;

    draw_editor(pane->pos_x + padding, pane->pos_y + padding, width, pane->split - padding);
    draw_preview(pane->pos_x + padding, pane->pos_y + pane->split + padding, width, height - pane->split - padding);
}
