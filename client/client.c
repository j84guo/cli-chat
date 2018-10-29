#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <pthread.h>
#include <arpa/inet.h>

#include "elp.h"

void stop_thread(pthread_t tid);
int arg_check(int argc, char **argv);
int stdin_loop(elpinfo_t *info);

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

int stdin_loop(elpinfo_t *info)
{
    char input[512];
    while (1) {
        if (!fgets(input, 512, stdin)) {
            if (ferror(stdin)) {
                perror("fgets");
                return 1;
            }

            break;
        }

        elp_msg_out(info, input);
    }

    return 0;
}

/**
 * todo:
 * - release mutex on cleanup
 * - does pthread_create need to be checked
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

    int ret = stdin_loop(&info);

    printf("Shutting down\n");
    stop_thread(elptid);
    elpinfo_dtry(&info);

    return ret;
}
