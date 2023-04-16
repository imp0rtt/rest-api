.PHONY:
.SILENT:

build:
	go build simpleServer/cmd/main

run: build
	go build simpleServer/cmd/main

