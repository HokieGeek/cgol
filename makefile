export GOPATH := ${PWD}

all: bin/cgol

bin/cgol: src/cgol.go src/cgol/pond.go src/cgol/inits.go src/cgol/rules.go
	go build -o $@ src/cgol.go 

clean:
	rm -rf bin

.PHONY: all clean
