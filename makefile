export GOPATH := ${PWD}

all: bin/cgol

bin/cgol: src/cgol.go src/cgol/universe.go
	go build -o $@ src/cgol.go 

clean:
	rm -rf bin

.PHONY: all clean
