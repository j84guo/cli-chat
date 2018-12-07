#include <errno.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>

#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>

#include "tcpcon.h"

int tcpcon_config(tcpcon_t *con, char *ip, unsigned short port)
{
    if (!con || !ip)
        return -1;

    if (!inet_pton(AF_INET, ip, &con->addr.sin_addr)) {
        fprintf(stderr, "tcpcon_init: invalid ip format\n");
        return -1;
    }

    con->addr.sin_family = AF_INET;
    con->addr.sin_port = htons(port);
    memset(&con->addr.sin_zero, 0, sizeof con->addr.sin_zero);

    if ((con->fd = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP)) == -1) {
        perror("socket");
        return -1;
    }

    return 0;
}

int tcpcon_dtry(tcpcon_t *con)
{
    if (!con)
        return -1;
    return close(con->fd);
}

int tcpcon_init(tcpcon_t *con, char *ip, unsigned short port)
{
    if (tcpcon_config(con, ip, port))
        return -1;

    if (connect(con->fd, (struct sockaddr *) &con->addr,
		sizeof con->addr) == -1) {
        perror("connect");
        return -1;
    }

    return 0;
}

int sendall(int fd, char *buf, int len)
{
    int i = 0, n;

    while (i < len) {
        if ((n = send(fd, buf+i, len-i, 0)) == -1) {
            if (errno == EINTR)
                continue;
            return -1;
        }
        i += n;
    }

    return 0;
}
