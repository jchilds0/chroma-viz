/*
 * editor.c 
 */

#include "chroma-viz.h" 

void draw_edit_rectangle(TILE *, SHOW *);
void draw_edit_panel(int, int, char *, int);
void rectangle_mouse_click(TILE *, Rect *);

const int padding = 30;

void draw_editor(TILE *editor, SHOW *show) {
    DrawRectangle(editor->pos_x, editor->pos_y, editor->width, editor->height, CHROMA_BG);
    DrawText("Editor", CENTER(editor->pos_x, editor->width), CENTER(editor->pos_y, editor->height), 20, CHROMA_TEXT);

    draw_edit_rectangle(editor, show);
}

void draw_edit_rectangle(TILE *editor, SHOW *show) {
    Graphic *graphic = &show->graphic[show->selected_page];
    char *attrs[4] = {"x: ", "y: ", "width: ", "height: "};
    int values[4] = {graphic->rect[0].pos_x, graphic->rect[0].pos_y, graphic->rect[0].width, graphic->rect[0].height};

    for (int i = 0; i < 4; i++) {
        draw_edit_panel(editor->pos_x + padding, editor->pos_y + padding + i * (padding + 50), attrs[i], values[i]);
    }
}

void draw_edit_panel(int pos_x, int pos_y, char *name, int value) {
    const int text_pad = 10;
    char value_str[20];
    memset(value_str, '\0', sizeof value_str);
    sprintf(value_str, "%d", value);

    DrawRectangle(pos_x, pos_y, 200, 50, WHITE);
    DrawText(name, pos_x + text_pad, pos_y + text_pad, 30, BLACK);

    DrawText(value_str, pos_x + (190 - MeasureText(value_str, 30)), pos_y + text_pad, 30, BLACK);

    DrawRectangle(pos_x + 210, pos_y, 50, 50, WHITE);
    DrawRectangle(pos_x + 270, pos_y, 50, 50, WHITE);

    DrawText("+", pos_x + 210 + text_pad, pos_y + text_pad / 2, 50, BLACK);
    DrawText("-", pos_x + 270 + text_pad, pos_y + text_pad / 2, 50, BLACK);
}

void editor_mouse_click(TILE *editor, SHOW *show) {
    //printf("Editor Click!\n");

    rectangle_mouse_click(editor, &show->graphic[show->selected_page].rect[0]);
}

void rectangle_mouse_click(TILE *editor, Rect *rectangle) {
    int *values_on[4] = {&rectangle->pos_x, &rectangle->pos_y, &rectangle->width, &rectangle->height};
    //int *values_off[4] = {&object_off->pos_x, &object_off->pos_y, &object_off->width, &object_off->height};

    int start_x, start_y;

    for (int i = 0; i < 4; i++) {
        start_x = editor->pos_x + 200 + padding;
        start_y = editor->pos_y + padding + i * (padding + 50);

        if (WITHIN(GetMouseX(), start_x, start_x + 50) && 
            WITHIN(GetMouseY(), start_y, start_y + 50)) {
            //printf("Click Value Up %d\n", i);

            (*values_on[i])++;

        }

        if (WITHIN(GetMouseX(), start_x + 50 + padding, start_x + 50 + padding + 50) &&
            WITHIN(GetMouseY(), start_y, start_y + 50)) {
            //printf("Click Value Down %d\n", i);

            (*values_on[i])--;
        }
    }
}
