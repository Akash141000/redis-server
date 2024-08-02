
build:
	@go build -o bin/redis

run: build
	./bin/redis

testPeer: 
	go run cmd/testPeer/testPeer.go