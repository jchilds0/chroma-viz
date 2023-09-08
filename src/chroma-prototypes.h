/*
 * chroma-prototypes.h 
 */

#ifndef CHROMA_CHROMA_PROTOTYPES
#define CHROMA_CHROMA_PROTOTYPES

#include <raylib.h>

/* chroma-output.c */
int connect_to_engine(char *, int);
int send_message_to_engine(int, char *);
int close_engine_connection(int);

/* editor.c */
void draw_editor(int, int, int, int);

/* preview.c */
void draw_preview(int, int, int, int);

/* templates.c */
void draw_templates(int, int, int, int);

/* show.c */
void draw_show(int, int, int, int);

#endif // !CHROMA_CHROMA_PROTOTYPES

