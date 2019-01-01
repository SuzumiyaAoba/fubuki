SRCS := \
	main.go \
	token/token.go \
	syntax/lexer.go

.PHONY: all fubuki

all: build

build: fubuki

fubuki: $(SRCS)
	go build;

run: build
	./fubuki

clean:
	$(RM) fubuki
	$(RM) y.output
	$(RM) -rf vendor
