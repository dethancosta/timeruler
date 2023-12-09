server:
	go build -o timeruler ./cmd/server

runServer:
	go run ./cmd/server

standalone:
	go run ./cmd/server -sa true &
