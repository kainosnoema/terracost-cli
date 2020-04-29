GO_FILES = $(find -iname '*.go')
BIN		 = ./terracost

$(BIN): $(GO_FILES)
	@echo "Building..."
	@CGO_ENABLED=0 go build -o $(BIN) main.go
	@echo "...done"

clean:
	@rm $(BIN)

test:
	@go test -v ./...

.PHONY: clean test
