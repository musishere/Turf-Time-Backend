server:
	nodemon --exec go run ./cmd/server/main.go --signal SIGTERM

build:
	go build -o bin/server cmd/server/main.go

run:
	./bin/server

clean:
	rm -f bin/server