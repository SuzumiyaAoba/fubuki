package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
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

	resID    = 0
	env      = make(lambda.Env)
	history  = filepath.Join(os.TempDir(), "fubuki.history")
	commands = map[string]func([]string){
		":load": loadFiles,
		":l":    loadFiles,
		":env":  showEnv,
		":help": help,
		":h":    help,
		":exit": exit,
		":show": show,
		":s":    show,
	}
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
		cmd := []string{":exit", ":help", ":env", ":load"}
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
			case strings.HasPrefix(line, ":exit"):
				goto exit
			case strings.HasPrefix(line, ":"): // command
				cmds := strings.Split(line, " ")
				exeCommand(cmds)
			case line == "": // ignore
				break
			default: // evaluate lambda expressions
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

func exeCommand(chunks []string) {
	cmd := chunks[0]
	options := make([]string, 0, len(chunks)-1)
	for _, c := range chunks[1:] {
		if c != "" {
			options = append([]string{c}, options...)
		}
	}
	reverseStrings(options)

	if exe, ok := commands[cmd]; ok {
		exe(options)
	} else {
		unknownCommand(cmd)
	}
}

func showEnv(options []string) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	if len(options) > 0 {
		switch options[0] {
		case "asc":
			sort.Strings(keys)
		case "desc":
			sort.Strings(keys)
			reverseStrings(keys)
		}
	}

	for _, key := range keys {
		fmt.Printf("%s := %s\n", key, lambda.Readable(env[key]))
	}
	fmt.Println()
}

func loadFiles(options []string) {
	for _, path := range options {
		s, err := locerr.NewSourceFromFile(path)
		if err != nil {
			fmt.Println(err)
			fmt.Println()
		} else {
			err := eval(s, true)
			if err == nil {
				green.Fprintf(os.Stdout, "Success: ")
				fmt.Printf("load %s\n", path)
			}
		}
	}
	fmt.Println()
}

func show(options []string) {
	for _, bind := range options {
		if term, ok := env[bind]; ok {
			green.Fprint(os.Stdout, "Exists: ")
			fmt.Printf("%s := %s\n", bind, lambda.Readable(term))
		} else {
			red.Fprint(os.Stdout, "Not found: ")
			fmt.Println(bind)
		}
	}
	fmt.Println()
}

func exit(options []string) {
	// Do nothing
}

func help(options []string) {
	// TODO
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
