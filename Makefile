all: server client-cur

vendor:
	glide install

server:
	go build srv/server.go

client-cur:
	go build current/client-cur.go

clean:
	-rm -f server
	-rm -f client-cur

.PHONY: vendor all clean server client-cur
