SRCS := \
	syntax/grammar.go \
	syntax/lexer.go \
	syntax/parser.go \
	token/token.go \
	ast/node.go \
	ast/visitor.go \
	ast/printer.go \
	lambda/id.go \
	lambda/alpha.go \
	lambda/beta.go \
	lambda/term.go \
	main.go

all: build

build: fubuki

fubuki: $(SRCS)
	go build;

syntax/grammar.go: syntax/grammar.go.y
	go get golang.org/x/tools/cmd/goyacc
	goyacc -o syntax/grammar.go syntax/grammar.go.y
