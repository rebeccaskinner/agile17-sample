all: server client-cur

vendor:
	glide install

server:
	go build srv/server.go

client-cur:
	go build current/client-cur.go

.PHONY: vendor all
