.PHONY: clean client server

client: client/client.c client/llist.c client/llist.h client/tcpcon.c client/tcpcon.h
	mkdir -p bin
	gcc -o bin/client -Wall -pthread -I client client/client.c client/llist.c client/tcpcon.c
	@echo "Built client into bin/client"

server: server/server.go
	mkdir -p bin
	go build -o bin/server server/server.go
	@echo "Built server into bin/server"

clean:
	rm -rf bin
