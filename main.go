package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/SuzumiyaAoba/fubuki/lambda"
	"github.com/SuzumiyaAoba/fubuki/syntax"
	"github.com/fatih/color"

	"github.com/rhysd/locerr"

	"github.com/peterh/liner"
)

var (
	bold  = color.New(color.Bold)
	red   = color.New(color.FgRed)
	blue  = color.New(color.FgBlue)
	green = color.New(color.FgGreen)

	version = "0.0.1"

	resID   = 0
	env     = make(lambda.Env)
	history = filepath.Join(os.TempDir(), "fubuki.history")
)

func binds() []string {
	bind := make([]string, 0, len(env))
	for k := range env {
		if !strings.HasPrefix(k, "#") {
			bind = append(bind, k)
		}
	}
	return bind
}

func complete(line string, pos int) (string, []string, string) {
	chunks := strings.Split(line, " ")
	chunk := chunks[len(chunks)-1]
	c := make([]string, 0)

	for _, n := range binds() {
		if strings.HasPrefix(n, chunk) {
			c = append(c, n)
		}
	}

	prefix := ""
	if len(chunks) > 1 {
		for i := 0; i < len(chunks)-1; i++ {
			prefix += chunks[i] + " "
		}
	} else {
		cmd := []string{":exit", ":help", ":env"}
		for _, n := range cmd {
			if strings.HasPrefix(n, chunk) {
				c = append(c, n)
			}
		}
	}

	return prefix, c, string([]rune(line)[pos:])
}

func main() {
	welcome()

	prompt := liner.NewLiner()
	defer prompt.Close()

	prompt.SetCtrlCAborts(true)
	prompt.SetTabCompletionStyle(liner.TabPrints)

	prompt.SetWordCompleter(complete)

	if f, err := os.Open(history); err == nil {
		prompt.ReadHistory(f)
		f.Close()
	}

	for {
		if line, err := prompt.Prompt("fubuki> "); err == nil {
			line = strings.Trim(line, " \t")
			if line != "" {
				prompt.AppendHistory(line)
			}
			switch {
			case strings.HasPrefix(line, ":"):
				cmds := strings.Split(line, " ")
				switch cmds[0] {
				case ":exit":
					goto exit
				case ":help", ":h":
					// TODO
				case ":env", ":e":
					showEnv()
				case ":load", ":l":
					loadFiles(cmds[1:])
				default:
					unknownCommand(cmds[0])
				}
			case line == "":
				break
			default:
				eval(&locerr.Source{
					Path:   "<stdin>",
					Code:   []byte(line),
					Exists: false,
				}, false)
			}
		} else if err == liner.ErrPromptAborted {
			goto exit
		} else if err == io.EOF {
			goto exit
		} else {
			fmt.Println(err)
		}
	}
exit:

	if f, err := os.Create(history); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		prompt.WriteHistory(f)
		f.Close()
	}
}

func welcome() {
	fmt.Printf("Welcome to Fubuki %s\n", version)
	fmt.Println("see https://github.com/SuzumiyaAoba/fubuki :help for help")
	fmt.Println()
}

func showEnv() {
	for k, v := range env {
		fmt.Printf("%s := %s\n", k, lambda.Readable(v))
	}
	fmt.Println()
}

func loadFiles(paths []string) {
	for _, path := range paths {
		s, err := locerr.NewSourceFromFile(path)
		if err != nil {
			fmt.Println(err)
			fmt.Println()
		} else {
			err := eval(s, true)
			if err == nil {
				green.Fprintf(os.Stdout, "success: ")
				fmt.Printf("load %s\n\n", path)
			}
		}
	}
}

func unknownCommand(cmd string) {
	red.Fprint(os.Stdout, "Error: ")
	bold.Fprintf(os.Stdout, "unknown command: %s\n\n", cmd)
}

func eval(source *locerr.Source, silent bool) error {
	t, err := syntax.Parse(source)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		terms := lambda.AstToTerms(t)
		alpha := lambda.Alpha(terms)
		beta := lambda.Beta(env, alpha)
		for _, term := range beta {
			if !silent {
				id := fmt.Sprintf("#%d", resID)
				fmt.Printf("%s: %s\n\n", id, lambda.Readable(term))
				env[id] = term
				resID++
			}
		}
	}
	return nil
}
