build:
	go build -o out/gocom main.go

container:
	docker build -t gocom .

debug-container:
	docker run -it --rm --name gocom -p 8080:8080 gocom

compose:
	docker-compose up -d

run-server:
	./out/gocom --mode server

run-local:
	go run main.go