SRCS := \
	cmd/fubuki/main.go \
	token/token.go \
	ast/node.go \
	ast/visitor.go \
	ast/printer.go \
	syntax/grammar.go \
	syntax/lexer.go \
	syntax/parser.go \
	lambda/id.go \
	lambda/alpha.go \
	lambda/beta.go \
	lambda/term.go \

.PHONY: all parser fubuki syntax/grammar.go

all: build

parser: syntax/grammar.go

build: fubuki

fubuki: $(SRCS)
	cd cmd/fubuki && go build;

syntax/grammar.go: syntax/grammar.go.y
	go get golang.org/x/tools/cmd/goyacc
	goyacc -o syntax/grammar.go syntax/grammar.go.y

clean:
	$(RM) syntax/grammar.go
	$(RM) fubuki
	$(RM) y.output
	$(RM) -rf vendor
