all:
	make build exec
build:
	go build -o ./bin ./cmd/tmail
exec:
	./bin/tmail

