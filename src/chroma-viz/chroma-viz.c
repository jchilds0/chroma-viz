/*
 * chroma-viz.c 
 */

#include "chroma-viz.h" 

#define PADDING     5

void update_window_panes(PANE *, PANE *, PANE *);
void update_pane_tiles(PANE *, TILE *, TILE *);

int main(int argc, char **argv) {
    PANE main  = {0, 0, 1200, 800, 600};
    PANE left, right;
    TILE editor, preview, templates, show_tile;
    bool resize_mid = false, resize_left = false;
    Connection engine = (Connection) {1920, 1080, ENGINE_DISCON, -1, "127.0.0.1", 6100};
    Connection prev = (Connection) {800, 450, ENGINE_DISCON, -1, "127.0.0.1", 6000};

    Keymap keymap = (Keymap) { KEY_KP_DIVIDE, KEY_KP_MULTIPLY, KEY_KP_ADD, KEY_KP_SUBTRACT };

    // Navbar
    main.pos_y = 20;
    main.height = main.height - main.pos_y - 30;
    left.split = main.height / 2;
    TILE lower = (TILE) {main.pos_x, main.pos_y + main.height, main.width, 30};
    //SHOW *show = init_show();
    SHOW *show = read_show_from_file("shows/basic_show.chromashow");
    //write_show_to_file(show);

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
            if (WITHIN(GetMouseY(), lower.pos_y, lower.pos_y + lower.height)) {

                lower_bar_mouse_click(&lower, &engine, &prev);

            } else if (WITHIN(GetMouseX(), show_tile.pos_x, show_tile.pos_x + show_tile.width) 
                && WITHIN(GetMouseY(), show_tile.pos_y, show_tile.pos_y + show_tile.height)) {

                show_mouse_click(&show_tile, show, &engine);

            } else if (WITHIN(GetMouseX(), editor.pos_x, editor.pos_x + editor.width)
                && WITHIN(GetMouseY(), editor.pos_y, editor.pos_y + editor.height)) {

                editor_mouse_click(&editor, show);

            } else if (WITHIN(GetMouseX(), templates.pos_x, templates.pos_x + templates.width)
                && WITHIN(GetMouseY(), templates.pos_y, templates.pos_y + templates.height)) {

                templates_mouse_click(&templates, show);

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

        handle_keypress(&keymap, show, &engine);

        if (resize_mid) {
            main.split = GetMouseX();
        } else if (resize_left) {
            left.split = GetMouseY() - left.pos_y;
        }

        draw_connection_button(&main, &lower, &engine, &prev);

        update_window_panes(&main, &left, &right);
        update_pane_tiles(&left, &templates, &show_tile);
        update_pane_tiles(&right, &editor, &preview);
        update_lower_bar(&main, &lower);

        draw_templates(&templates);
        draw_show(&show_tile, show);
        draw_editor(&editor, show);
        draw_preview(&preview, &engine, &show->graphic[show->selected_page]);

        EndDrawing();
    }

    if (engine.status == ENGINE_CON) {
        close_engine_connection(engine.socket_desc);
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

