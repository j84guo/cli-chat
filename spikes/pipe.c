#include <stdio.h>
#include <unistd.h>
#include <pthread.h>
#include <errno.h>

void *pipe_thread(void *arg)
{
    int fd = *(int *) arg, ret;
    char buf;

    while ((ret = read(fd, &buf, 1))) {
        if (ret == -1) {
            if (errno == EINTR)
                continue;

            perror("read");
            break;
        }

        printf("%c", buf);
    }

    return NULL;
}

int main()
{
    pthread_t t1;
    int pfd[2];

    if (pipe(pfd) == -1) {
        perror("pipe");
        return 1;
    }

    pthread_create(&t1, NULL, pipe_thread, &pfd[0]);
    
    char buf;
    int ret;
    while ((ret = read(STDIN_FILENO, &buf, 1))) {
        if (ret == -1) {
            if (errno == EINTR)                                                 
                 continue;                                                       
                                                                                 
            perror("read");                                                     
            return 1;
        }

        while((ret = write(pfd[1], &buf, 1) != 1)) {
            if (ret == -1) {
                if (errno == EINTR)
                    continue;

                perror("write");
                return 1;
            }
        }
    }

    return 0;
}
