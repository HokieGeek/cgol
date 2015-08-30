export GOPATH := ${PWD}

all: bin/cgol-cli bin/cgold

bin/cgol-cli: src/cgol-cli.go src/cgol/pond.go src/cgol/inits.go src/cgol/rules.go src/cgol/strategy.go src/cgol/processors.go src/cgol/gameboard.go
	go build -o $@ src/cgol-cli.go 

bin/cgold: src/cgol-cli.go src/cgol/pond.go src/cgol/inits.go src/cgol/rules.go src/cgol/strategy.go src/cgol/processors.go src/cgol/gameboard.go
	go build -o $@ src/cgol-cli.go 

run:
	go run src/cgol-cli.go

test:
	@echo "[Running unit tests]"
	@go test cgol

clean:
	go clean
	rm -rf bin

.PHONY: all run test clean
