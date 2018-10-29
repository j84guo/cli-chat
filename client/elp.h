#include "llist.h"
#include "tcpcon.h"

typedef struct {
    int pipefd[2];
    tcpcon_t con;
    llist_t queue;
    pthread_mutex_t mutex;
} elpinfo_t;

int elpinfo_init(elpinfo_t *info, char *ip, unsigned short port);

int elpinfo_dtry(elpinfo_t *info);

void *elp_run(void *arg);

int elp_loop(elpinfo_t *info);

int elp_give_msg(elpinfo_t *info, char *input);

int fdset_init(fd_set *set, int pd, int sd);
