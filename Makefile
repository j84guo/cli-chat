chat:
	go build -o bin/server cmd/server.go
	go build -o bin/client cmd/client.go

clean:
	rm -rf bin

.PHONY: chat
