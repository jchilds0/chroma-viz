/*
 * chroma-prototypes.h 
 */

#ifndef CHROMA_CHROMA_PROTOTYPES
#define CHROMA_CHROMA_PROTOTYPES

#include <raylib.h>

/* chroma-output.c */
void open_socket_connection(void);

/* editor.c */
void draw_editor(int, int, int, int);

/* preview.c */
void draw_preview(int, int, int, int);

/* templates.c */
void draw_templates(int, int, int, int);

/* show.c */
void draw_show(int, int, int, int);

#endif // !CHROMA_CHROMA_PROTOTYPES

