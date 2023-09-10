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
int render_objects(int, RenderObject *, int);

/* editor.c */
void draw_editor(TILE *);

/* preview.c */
void draw_preview(TILE *);

/* templates.c */
void draw_templates(TILE *);

/* show.c */
void draw_show(TILE *, SHOW *);
SHOW *init_show(void);
void free_show(SHOW *);
void show_mouse_click(TILE *, SHOW *, int, int);

#endif // !CHROMA_CHROMA_PROTOTYPES

