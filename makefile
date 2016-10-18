package = github.com/WianVos/goxld

.PHONY: install release test

install:
	go get -t -v ./...

release:
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o release/goxld-linux-amd64 $(package)
	GOOS=linux GOARCH=386 go build -o release/goxld-linux-386 $(package)
	GOOS=linux GOARCH=arm go build -o release/goxld-linux-arm $(package)

test:
	go test -v
