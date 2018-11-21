client/client: bin/llist.o bin/tcpcon.o bin/elp.o bin/client.o
	mkdir -p bin
	gcc -o bin/client bin/client.o bin/llist.o bin/tcpcon.o bin/elp.o -pthread
	echo "Built client into bin/client"

bin/llist.o:
	mkdir -p bin
	gcc -c -o bin/llist.o -Wall -pthread -I client client/llist.c

bin/tcpcon.o:
	mkdir -p bin
	gcc -c -o bin/tcpcon.o -Wall -pthread -I client client/tcpcon.c

bin/elp.o:
	mkdir -p bin
	gcc -c -o bin/elp.o -Wall -pthread -I client client/elp.c

bin/client.o:
	mkdir -p bin
	gcc -c -o bin/client.o -Wall -pthread -I client client/client.c

server: client/client
	mkdir -p bin
	go build -o bin/server server/server.go
	echo "Built server into bin/server"

