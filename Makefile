chat:
	go build -o bin/server server.go
	go build -o bin/client client.go

clean:
	rm -rf bin

.PHONY: chat
