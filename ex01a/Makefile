go_files := $(wildcard *.go)
tags := netgo,osusergo
ldflags := "-extldflags '-static' -s -w "
executable := quiz

runQuiz: build
	./quiz

print:
	echo $(go_files)

build: print
	CGO_ENABLED=0 go build -a -tags $(tags) -ldflags $(ldflags) -o $(executable) $(go_files)
