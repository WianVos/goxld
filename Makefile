build_all: xld.build goxld.build

goxld.build:
	go build ./
	go install ./

xld.build:
	go build ../xld
	go install ../xld
