#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <pthread.h>

#include "elp.h"

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

    pthread_mutex_init(&info->mutex, NULL);

    return 0;
}

int elpinfo_dtry(elpinfo_t *info)
{
    int ret = 0;

    if (llist_dtry(&info->queue))
        ret = -1;

    if (tcpcon_dtry(&info->con))
        ret = -1;

    if (close(info->pipefd[0]))
        ret = -1;

    if (close(info->pipefd[1]))
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

int elp_loop(elpinfo_t *info)
{
    int ret, nfds, pd=info->pipefd[0], sd=info->con.fd;
    char rbuf[512 + 1], *msg;
    fd_set set;

    while (1) {
        fdset_init(&set, pd, sd);
        nfds = pd >= sd ? pd + 1 : sd + 1;
        ret = select(nfds, &set, NULL, NULL, NULL);

        if (ret == -1) {
            if (errno == EINTR)
                continue;

            perror("select");
            return -1;
        }

        if (FD_ISSET(pd, &set)) {
            read(pd, rbuf, 1);

            pthread_mutex_lock(&info->mutex);
            msg = llist_remf(&info->queue);
            pthread_mutex_unlock(&info->mutex);

            sendall(sd, msg, strlen(msg));
            free(msg);
        } else {
            ret = recv(sd, rbuf, 512, 0);

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
    elp_loop(info);
    return NULL;
}

int elp_give_msg(elpinfo_t *info, char *input)
{
    if (!info || !input)
        return -1;

    int len = strlen(input) + 1;
    char *msg = malloc(len);
    strncpy(msg, input, len);

    pthread_mutex_lock(&info->mutex);
    llist_addl(&info->queue, msg);
    pthread_mutex_unlock(&info->mutex);

    write(info->pipefd[1], "1", 1);
    return 0;
}
