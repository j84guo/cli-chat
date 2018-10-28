client: client.c llist.c llist.h tcpcon.c tcpcon.h
	gcc -o client -Wall -pthread client.c llist.c tcpcon.c
