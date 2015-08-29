export GOPATH := ${PWD}

all: bin/cgol

bin/cgol: src/cgol.go src/cgol/pond.go src/cgol/inits.go src/cgol/rules.go src/cgol/strategy.go src/cgol/processors.go src/cgol/gameboard.go
	go build -o $@ src/cgol.go 

test:
	@echo "[Running unit tests]"
	@go test cgol

clean:
	rm -rf bin

.PHONY: all clean
