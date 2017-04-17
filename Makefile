all: server client-cur client-monadic

vendor:
	glide install

server:
	go build srv/server.go

client-cur:
	go build current/client-cur.go

client-monadic:
	go build monadic/client-monadic.go

clean:
	-rm -f server
	-rm -f client-cur
	-rm -f client-monadic

.PHONY: vendor all clean server client-cur client-monadic
