/*
 * chroma-prototypes.h 
 */

#ifndef CHROMA_CHROMA_PROTOTYPES
#define CHROMA_CHROMA_PROTOTYPES

#include <raylib.h>
#include "chroma-typedefs.h"

/* chroma-output.c */
int connect_to_engine(char *, int);
int send_message_to_engine(int, char *);
int close_engine_connection(int);

/* chroma-renderer.c */ 
void render_graphic(Graphic *, Connection *, bool);

/* editor.c */
void draw_editor(TILE *, SHOW *);
void editor_mouse_click(TILE *, SHOW *);

/* preview.c */
void draw_preview(TILE *, Connection *, Graphic *);

/* templates.c */
void draw_templates(TILE *);
void templates_mouse_click(TILE *, SHOW *);

/* show.c */
void draw_show(TILE *, SHOW *);
SHOW *init_show(void);
void free_show(SHOW *);
void show_mouse_click(TILE *, SHOW *, Connection *);
void animate_on(SHOW *, Connection *, int);
void animate_off(SHOW *, Connection *, int);

/* show_io.c */ 
void write_show_to_file(SHOW *);
SHOW *read_show_from_file(char *);

/* lower_bar.c */ 
void draw_connection_button(PANE *, TILE *, Connection *, Connection *);
void update_lower_bar(PANE *, TILE *);
void lower_bar_mouse_click(TILE *, Connection *, Connection *);

/* keymap.c */ 
void handle_keypress(Keymap *, SHOW *, Connection *);

#endif // !CHROMA_CHROMA_PROTOTYPES

