/*
 * chroma-output.c 
 */

#include "chroma-viz.h"

#include <stdio.h>
#include <string.h>
#include <sys/socket.h>
#include <arpa/inet.h>

int connect_to_engine(char *addr, int port) {
    int socket_desc;
    struct sockaddr_in server_addr;

    // create socket 
    socket_desc = socket(AF_INET, SOCK_STREAM, 0);

    if (socket_desc < 0) {
        printf("Unable to create socket\n");
        return -1;
    }

    printf("Socket created successfully\n");

    // set port and ip
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(6100);
    server_addr.sin_addr.s_addr = inet_addr("127.0.0.1");

    // send connection request to server
    if (connect(socket_desc, (struct sockaddr*)&server_addr, sizeof server_addr) < 0) {
        printf("Unable to connect\n");
        return -1;
    }

    printf("Connected with server successfully\n");
    return socket_desc;
}

int send_message_to_engine(int socket_desc, char *message) {
    char server_message[2000];

    // clean buffers
    memset(server_message, '\0', sizeof server_message );

    // send message to server 
    if (send(socket_desc, message, strlen(message), 0) < 0) {
        printf("Unable to send message\n");
        return -1;
    }

    // recieve the servers response
    if (recv(socket_desc, server_message, sizeof server_message, 0) < 0) {
        printf("Error while recieving server's msg\n");
        return -1;
    }

    //printf("Server's response: %s\n", server_message);

    return 0;
}

int close_engine_connection(int socket_desc) {
    char msg[1];
    msg[0] = 4;

    if (send_message_to_engine(socket_desc, msg) < 0) {
        return -1;
    }
    shutdown(socket_desc, SHUT_RDWR);

    return 1;
}
