cli:
	go build -o tc-cli ./cmd/cli

tui:
	go build -o tc-tui ./cmd/tui

server:
	go build -o timeruler ./cmd/server
