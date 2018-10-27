client: client.c llist.c llist.h
	gcc -o client -Wall -pthread client.c llist.c
