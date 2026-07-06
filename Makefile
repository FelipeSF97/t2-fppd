EXAMPLES = \
01-hello \
02-ping-pong \
03-ring \
04-scatter-gather \
05-reduce-pi \
06-workpool \
t2/paralelo

.PHONY: all clean deps

all: deps $(patsubst %,bin/%,$(EXAMPLES))

deps:
	@go mod tidy

bin/%: %/main.go
	@echo "Compilando $*..."
	@mkdir -p bin
	@go build -o $@ ./$*

clean:
	rm -rf bin/
