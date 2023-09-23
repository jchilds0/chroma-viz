/*
 * chroma-engine.h 
 */

#ifndef CHROMA_CHROMA_ENGINE
#define CHROMA_CHROMA_ENGINE

#include <raylib.h>
#include "tcp.h"

#define CHROMA_FRAMERATE              30

#define CHROMA_TIMEOUT                -1
#define CHROMA_CLOSE_SOCKET           -2

int start_tcp_server(char *, int);
int listen_for_client(int);
int recieve_message(int, char *);

#endif // !CHROMA_CHROMA_ENGINE
