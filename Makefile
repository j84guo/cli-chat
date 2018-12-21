chat:
	go build -o bin/server server.go utils.go
	go build -o bin/client client.go utils.go config.go

clean:
	rm -rf bin

.PHONY: chat
