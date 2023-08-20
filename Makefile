build:
	go build -o redis-lite .

run: build
	./redis-lite

dev:
	go run main.go