/*
 * keymap.c 
 */

#include "chroma-viz.h"
#include <stdbool.h>

void handle_keypress(Keymap *keymap, SHOW *show, Connection *conn) {
    int key = GetKeyPressed();

    if (key == keymap->animate_on) {

        graphic_to_engine(conn->socket_desc, &show->graphic[show->selected_page], true);

    } else if (key == keymap->animate_off) {

        graphic_to_engine(conn->socket_desc, &show->graphic[show->selected_page], false);

    } else if (key == keymap->next_page) {

        if (show->selected_page < show->num_pages - 1)
            show->selected_page++;

    } else if (key == keymap->prev_page) {

        if (show->selected_page > 0)
            show->selected_page--;

    }

}
