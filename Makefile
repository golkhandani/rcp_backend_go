hello:
	echo "Hello"

build:
	go build -o bin/main ./main.go

run:
	go run ./main.go

run-dev:
	npx nodemon --exec "go run" ./main.go --signal SIGTERM
