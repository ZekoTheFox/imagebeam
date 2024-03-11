.PHONY: run clean-win build

run:
	go run cmd/imagebeam-server.go
 
# expects `DISCORD_TOKEN` env to be set 
run-dev: build
#	go build -o bin/ cmd/imagebeam-server.go 
	bin/imagebeam-server.exe -token ${DISCORD_TOKEN}

clean-win:
	del bin/imagebeam-server.exe

build:
	go build -o bin/ cmd/imagebeam-server.go
