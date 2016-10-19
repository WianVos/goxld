package = github.com/WianVos/goxld

.PHONY: install release test

install:
	go get -t -v ./...

release:
	mkdir -p release
	rm -rf release/*
	GOOS=linux GOARCH=amd64 go build -o release/goxld-linux-amd64 $(package)
	GOOS=linux GOARCH=386 go build -o release/goxld-linux-386 $(package)
	GOOS=windows GOARCH=amd64 go build -o release/goxld-windows-amd64 $(package)
	GOOS=windows GOARCH=386 go build -o release/goxld-windows-386 $(package)
	GOOS=darwin GOARCH=amd64 go build -o release/goxld-osx-amd64 $(package)
test:
	go test -v
