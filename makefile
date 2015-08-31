export GOPATH := ${PWD}

all: bin/cgol-cli bin/cgold

bin/cgol-cli: src/cgol/cgol-cli.go src/cgol/core/pond.go src/cgol/core/patterns.go src/cgol/core/rules.go src/cgol/core/strategy.go src/cgol/core/processors.go src/cgol/core/gameboard.go
	go build -o $@ src/cgol/cgol-cli.go 

bin/cgold: src/cgol/cgol-cli.go src/cgol/core/pond.go src/cgol/core/patterns.go src/cgol/core/rules.go src/cgol/core/strategy.go src/cgol/core/processors.go src/cgol/core/gameboard.go
	go build -o $@ src/cgol/cgol-cli.go 

run:
	go run src/cgol/cgol-cli.go

test:
	@echo "[Running unit tests]"
	@go test cgol/core

clean:
	go clean
	rm -rf bin

.PHONY: all run test clean
