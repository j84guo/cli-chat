#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <pthread.h>
#include <arpa/inet.h>

#include "llist.h"
#include "tcpcon.h"

typedef struct {
    int pipefd[2];
    tcpcon_t con;
    llist_t queue;
} elpinfo_t;

int elpinfo_init(elpinfo_t *info, char *ip, unsigned short port);
int elpinfo_destroy(elpinfo_t *info);
void *elp_run(void *arg);
int elp_loop(int pd, tcpcon_t *con, llist_t *queue);
int elp_msg(elpinfo_t *info, char *input);
void stop_thread(pthread_t tid);
int arg_check(int argc, char **argv);
int fdset_init(fd_set *set, int pd, int sd);

int elpinfo_init(elpinfo_t *info, char *ip, unsigned short port)
{
    if (!info || !ip)
        return -1;

    llist_init(&info->queue);

    if (pipe(info->pipefd) == -1) {
        perror("pipe");
        return -1;
    }

    if (tcpcon_init(&info->con, ip, port))
        return -1;

    return 0;
}

int elpinfo_destroy(elpinfo_t *info)
{
    int ret = 0;

    if (llist_destroy(&info->queue))
        ret = -1;

    if (tcpcon_destroy(&info->con))
        ret = -1;

    return ret;
}

int fdset_init(fd_set *set, int pd, int sd)
{
    if (!set || pd < 0 || sd < 0)
        return -1;

    FD_ZERO(set);
    FD_SET(pd, set);
    FD_SET(sd, set);

    return 0;
}

int elp_loop(int pd, tcpcon_t *con, llist_t *queue)
{
    int ret, nfds;
    char rbuf[512 + 1], *msg;
    fd_set readfds;

    while (1) {
        fdset_init(&readfds, pd, con->fd);
        nfds = pd >= con->fd ? pd + 1 : con->fd + 1;
        ret = select(nfds, &readfds, NULL, NULL, NULL);

        if (ret == -1) {
            if (errno == EINTR)
                continue;

            perror("select");
            return -1;
        }

        if (FD_ISSET(pd, &readfds)) {
            read(pd, rbuf, 1);

            msg = llist_remf(queue);
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

void *elp_run(void *arg)
{
    elpinfo_t *info = (elpinfo_t *) arg;
    elp_loop(info->pipefd[0], &info->con, &info->queue);
    return NULL;
}

void stop_thread(pthread_t tid)
{
    pthread_cancel(tid);
    pthread_join(tid, NULL);
}

int arg_check(int argc, char **argv)
{
    if (argc != 3) {
        fprintf(stderr, "Usage: %s <ip> <port>\n", argv[0]);
        return -1;
    } else if (strlen(argv[1]) > INET6_ADDRSTRLEN) {
        fprintf(stderr, "Usage: <ip> is over %d bytes\n", INET6_ADDRSTRLEN);
        return -1;
    }

    return 0;
}

int elp_msg(elpinfo_t *info, char *input)
{
    if (!info || !input)
        return -1;

    int len = strlen(input) + 1;
    char *msg = malloc(len);
    strncpy(msg, input, len);
    llist_addl(&info->queue, msg);
    write(info->pipefd[1], "1", 1);

    return 0;
}

/**
 * todo:
 * - lock queue access
 * - cancel thread with mutex
 */
int main(int argc, char **argv)
{
    if (arg_check(argc, argv))
        return 1;

    printf("Starting up\n");
    elpinfo_t info;
    if (elpinfo_init(&info, argv[1], atoi(argv[2])))
        return 1;

    pthread_t elptid;
    pthread_create(&elptid, NULL, elp_run, &info);

    char input[512];
    while (1) {
        if (!fgets(input, 512, stdin)) {
            if (ferror(stdin)) {
                perror("fgets");
                stop_thread(elptid);
                return 1;
            }

            break;
        }

        elp_msg(&info, input);
    }

    printf("Shutting down\n");
    stop_thread(elptid);
    elpinfo_destroy(&info);

    return 0;
}
