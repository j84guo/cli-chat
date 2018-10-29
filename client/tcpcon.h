#include <netinet/in.h>

typedef struct {
    int fd;
    struct sockaddr_in addr;
} tcpcon_t;

int tcpcon_config(tcpcon_t *con, char *ip, unsigned short port);

int tcpcon_dtry(tcpcon_t *con);

int tcpcon_init(tcpcon_t *con, char *ip, unsigned short port);

int sendall(int fd, char *buf, int len);
