/*
 * keymap.c 
 */

#include "chroma-prototypes.h"
#include "chroma-viz.h"

void handle_keypress(Keymap *keymap, SHOW *show, Connection *conn) {
    int key = GetKeyPressed();

    if (key == keymap->animate_on) {

        render_graphic(&show->graphic[show->selected_page], conn, true);

    } else if (key == keymap->animate_off) {

        render_graphic(&show->graphic[show->selected_page], conn, false);

    } else if (key == keymap->next_page) {

        if (show->selected_page < show->num_pages - 1)
            show->selected_page++;

    } else if (key == keymap->prev_page) {

        if (show->selected_page > 0)
            show->selected_page--;

    }

}
