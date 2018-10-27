#include <stdio.h>
#include <unistd.h>
#include <sys/time.h>
#include <sys/types.h>

#define TIMEOUT 5
#define BUFLEN 1024

int main()
{
    struct timeval tv;
    fd_set readset;
    int ret;

    FD_ZERO(&readset);
    FD_SET(STDIN_FILENO, &readset);

    tv.tv_sec = TIMEOUT;
    tv.tv_usec = 0;

    ret = select(
        STDIN_FILENO + 1,
        &readset,
        NULL,
        NULL,
        &tv
    );

    if (ret == -1) {
        perror("select");
        return 1;
    } else if (ret == 0) {
        printf("%d seconds are up\n", TIMEOUT);
        return 0;
    }

    /**
     * must be true since 0 was the only file descriptor passed in
     */   
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
