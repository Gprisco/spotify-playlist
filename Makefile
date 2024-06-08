test:
	go test ./...

run:
	go run main.go

build:
	go build

benchmark:
	go test -v ./... -bench=. -run=xxx -benchmem
