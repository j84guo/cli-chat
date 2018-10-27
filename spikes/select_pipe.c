#include <stdio.h>
#include <unistd.h>
#include <sys/time.h>
#include <sys/types.h>

#define BUFLEN 1024

int main()
{
    int pfd[2];
    if (pipe(pfd) == -1) {
        perror("pipe");
        return 1;
    }

    fd_set readset;
    int ret;

    FD_ZERO(&readset);
    FD_SET(STDIN_FILENO, &readset);
    FD_SET(pfd[0], &readset);

    ret = select(pfd[0] + 1, &readset, NULL, NULL, NULL);
    if (ret == -1) {
        perror("select");
        return 1;
    }

    if (FD_ISSET(STDIN_FILENO, &readset)) {
        char buf[BUFLEN + 1];
        int len = read(STDIN_FILENO, buf, BUFLEN);

        if (len == -1) {
            perror("read");
            return 1;
        }

        if (len) {
            buf[len] = '\0';
            printf("read: %s", buf);
        }
    }

    return 0;
}
