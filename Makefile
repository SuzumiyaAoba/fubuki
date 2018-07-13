SRCS := \
	main.go \
	token/token.go \
	ast/node.go \
	ast/visitor.go \
	ast/printer.go \
	syntax/lexer.go \
	syntax/grammar.go \
	syntax/parser.go

all: build

build: fubuki

fubuki: $(SRCS)
	go build;

syntax/grammar.go: syntax/grammar.go.y
	go get golang.org/x/tools/cmd/goyacc
	goyacc -o syntax/grammar.go syntax/grammar.go.y
