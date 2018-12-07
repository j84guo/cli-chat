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
    if (pthread_mutex_destroy(&info->mutex))
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

void elp_cleanup(void *arg)
{
    pthread_mutex_unlock(arg);
}

void elp_msg_send(elpinfo_t *info)
{
    char c, *msg;
    read(info->pipefd[0], &c, 1);
    pthread_cleanup_push(elp_cleanup, &info->mutex);
    pthread_mutex_lock(&info->mutex);

	/* critical section */
    msg = llist_remf(&info->queue);

    pthread_cleanup_pop(1);
    sendall(info->con.fd, msg, strlen(msg));
    free(msg);
}

int elp_select(fd_set *set, int pd, int sd)
{
    if(fdset_init(set, pd, sd))
        return -1;
     
    return select(
        pd >= sd ? pd + 1 : sd + 1, 
        set, 
        NULL, 
        NULL, 
        NULL
    );
}

/**
 * Returns -1 to indicate recv error, 0 for server disconnect, 1 for success.
 * This error code convention is different from the other functions (0 on
 * success, -1 on failure) so maybe this code should be changed later...
 */
int elp_msg_recv(int sd)
{
    char rbuf[512 + 1];
    int ret = recv(sd, rbuf, 512, 0);

    if (ret == -1) {
        perror("recv");
        return -1;
    } else if (!ret) {
        printf("Server disconected\n");
        return 0;
    }

    rbuf[ret] = '\0';
    printf("%s", rbuf);
    return 1;
}

int elp_loop(elpinfo_t *info)
{
    int ret, pd=info->pipefd[0], sd=info->con.fd;
    fd_set set;

    while (1) {
        ret = elp_select(&set, pd, sd);   
    
        if (ret == -1) {
            if (errno == EINTR)
                continue;

            perror("elp_select");
            return -1;
        }

        if (FD_ISSET(pd, &set)) {
            elp_msg_send(info);
        } else {
            ret = elp_msg_recv(sd);

            if (ret == -1)
                return ret;
            else if (!ret)
                break;
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

int elp_msg_out(elpinfo_t *info, char *input)
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
