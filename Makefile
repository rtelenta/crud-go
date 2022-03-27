# LOCAL

install:
	go mod download

dev:
	gin -p 3001 -a 3007 -i run main.go
