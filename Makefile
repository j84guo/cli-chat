bin/client: client/client.c client/llist.c client/llist.h client/tcpcon.c client/tcpcon.h
	mkdir -p bin
	gcc -o bin/client -Wall -pthread -I client client/client.c client/llist.c client/tcpcon.c
