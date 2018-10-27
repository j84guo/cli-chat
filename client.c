#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <pthread.h>

#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>

#include "llist.h"

typedef struct {
    int fd;
    struct sockaddr_in addr;
} tcpcon_t;

typedef struct {
    int pd;
    tcpcon_t con;
    llist_t *queue;
} elpinfo_t;

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

int tcpcon_destroy(tcpcon_t *con)
{
    if (!con)
        return -1;

    return close(con->fd);
}

int tcpcon_init(tcpcon_t *con, char *ip, unsigned short port)
{
    int res = tcpcon_config(con, ip, port);
    if (res != 0)
        return res;

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
        if ((n = send(fd, buf, len, 0)) == -1) {
            if (errno == EINTR)
                continue;

            return -1;
        }

        i += n;
    }

    return 0;
}

int loop(int pd, tcpcon_t *con, llist_t *queue)
{
    int ret;
    char rbuf[512 + 1];
    fd_set readfds;

    while (1) {
        FD_ZERO(&readfds);
        FD_SET(pd, &readfds);
        FD_SET(con->fd, &readfds);
        
        ret = select(
            pd >= con->fd ? pd + 1 : con->fd + 1,
            &readfds, NULL, NULL, NULL);
        if (ret == -1) {
            if (errno == EINTR)
                continue;

            perror("select");
            return -1;
        }

        if (FD_ISSET(pd, &readfds)) {
            read(pd, rbuf, 1);

            char *msg = llist_remf(queue);
            sendall(con->fd, msg, strlen(msg));
            free(msg);
        } else {
            ret = recv(con->fd, rbuf, 512, 0);

            if (ret == -1) {
                perror("recv");
                return -1;
            } else if (!ret) {
                printf("Server disconected\n");
                return 0;
            }

            rbuf[ret] = '\0';
            printf("%s", rbuf);
        }

    }

    return 0;
}

void *run_eloop(void *arg)
{
    elpinfo_t *info = (elpinfo_t *) arg;
    loop(info->pd, &info->con, info->queue);
    return NULL;
}

void stop_thread(pthread_t tid)
{
    pthread_cancel(tid);
    pthread_join(tid, NULL);
}

int main(int argc, char **argv)
{
    if (argc != 3) {
        fprintf(stderr, "Usage: %s <ip> <port>\n", argv[0]);
        return 1;
    } else if (strlen(argv[1]) > INET6_ADDRSTRLEN) {
        fprintf(stderr, "Usage: <ip> less then  %d bytes\n", INET6_ADDRSTRLEN);
        return 1;
    }

    struct llist_t queue;
    llist_init(&queue);
    
    int pipefd[2];
    if (pipe(pipefd) == -1) {
        perror("pipe");
        return 1;
    }

    printf("Opening connection\n"); 
    elpinfo_t einfo;
    einfo.pd = pipefd[0];
    einfo.queue = &queue;
    if (tcpcon_init(&einfo.con, argv[1], atoi(argv[2])))
        return 1;

    pthread_t eloop;
    pthread_create(&eloop, NULL, run_eloop, &einfo);

    char input[512], *msg;
    int len;
    while (1) {
        if (!fgets(input, 512, stdin)) {
            if (ferror(stdin)) {
                perror("fgets");
                stop_thread(eloop);
                return 1;
            }

            break;
        }

        len = strlen(input) + 1;
        msg = malloc(len);
        strncpy(msg, input, len);
        llist_addl(&queue, msg);
        write(pipefd[1], "1", 1);
    }

    printf("Closing connection\n");
    stop_thread(eloop);
    tcpcon_destroy(&einfo.con);
    return 0;
}
