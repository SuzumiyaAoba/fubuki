package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/SuzumiyaAoba/fubuki/lambda"
	"github.com/SuzumiyaAoba/fubuki/syntax"

	"github.com/rhysd/locerr"

	"github.com/chzyer/readline"
)

func listFiles(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem(":load",
		readline.PcItemDynamic(listFiles("./")),
	),
	readline.PcItem(":exit"),
)

func main() {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[34mfubuki>\033[0m ",
		HistoryFile:     "/tmp/fubuki.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	env := make(lambda.Env)
	resID := 0
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case line == ":exit":
			goto exit
		case line == ":help":
		case line == "":
		default:
			t, _ := syntax.Parse(&locerr.Source{
				"<stdin>",
				[]byte(line),
				false,
			})
			terms := lambda.AstToTerms(t)
			alpha := lambda.Alpha(terms)
			beta := lambda.Beta(env, alpha)
			for _, term := range beta {
				id := fmt.Sprintf("#%d", resID)
				fmt.Printf("%s: %s\n\n", id, lambda.Readable(term))
				env[id] = term
				resID++
			}
		}
	}
exit:
}
