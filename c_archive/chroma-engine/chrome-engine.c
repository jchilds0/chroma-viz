/*
 * chroma-engine.c 
 */

#include "chroma-engine.h"
#include <sys/socket.h>

void render_text(char *);
void render_graphics(Graphic *);
void write_pixel_str_to_pixels(Color *, char *);


int main(int argc, char **argv) {
    const int screen_width = 1280;
    const int screen_height = 720;

    InitWindow(screen_width, screen_height, "raylib [core] example - basic window");
    SetTargetFPS(CHROMA_FRAMERATE);

    int socket_engine = start_tcp_server("127.0.0.1", 6100);
    int socket_client = -1, rec;
    Graphic renderer;

    char buf[MAX_BUF_SIZE];
    memset(buf, '\0', sizeof buf );

    renderer.num_rect = 0;
    renderer.num_text = 0;
    
    while (!WindowShouldClose()) {
        BeginDrawing();

        if (socket_client < 0) {
            socket_client = listen_for_client(socket_engine);
        } else {
            while ((rec = recieve_message(socket_client, buf)) == END_OF_GRAPHICS) {
                string_to_graphic(&renderer, buf);
                printf("%s\n", buf);
                memset(buf, '\0', sizeof buf);
            }

            if (rec == CHROMA_CLOSE_SOCKET) {
                shutdown(socket_client, SHUT_RDWR);
                socket_client = -1;
            } 
        }
        render_graphics(&renderer);

        EndDrawing();
    }

    shutdown(socket_engine, SHUT_RDWR);

    CloseWindow();
    return 0;
}

void render_text(char *buf) {
    DrawText(buf, 190, 200, 20, RAYWHITE);
}

void render_graphics(Graphic *render) {
    Rect *rect;
    Text *text;

    // render rectangles
    for (int j = 0; j < render->num_rect; j++) {
        rect = &render->rect[j];

        DrawRectangle(rect->pos_x, rect->pos_y, rect->width, rect->height, rect->color);
    }

    // render text 
    for (int j = 0; j < render->num_text; j++) {
        text = &render->text[j];

        DrawText(text->text, text->pos_x, text->pos_y, text->font_size, text->color);
    }
}

