package repl

import (
	"bufio"
	"fmt"
	"os"
	"pmp/lexml"
	"strings"
)

// Repl Type stores the Commands of the Repl, the FullSet that has been Parsed with lexml and the Values that the user process
type Repl struct {
	Commands map[string]func(repl Repl, index int) (Repl, error)
	FullSet  *lexml.Set
	Values   map[string]lexml.Data
	Prompt   string
	Keywords []string
}

func (repl Repl) ParsePrompt() ([]Command, error) {

	repl.Prompt = strings.Trim(repl.Prompt, "\n")
	lexer := NewLexer([]byte(repl.Prompt))
	tokenizer := NewTokenizer(repl, &lexer)
	tokenizer.Tokenize()

	return []Command{}, nil
}

// to replace by another thing
func (repl Repl) keywordCheck(word string) (bool, uint8) {
	if _, ok := repl.Commands[word]; ok {
		return true, COMMAND

	} else {
		for _, keyword := range repl.Keywords {
			if word == keyword {
				return true, KEYWORD
			}
		}
		return false, NONE
	}
}

func Launch(set *lexml.Set) error {

	var err error
	reader := bufio.NewReader(os.Stdin)
	repl := Repl{
		Commands: map[string]func(repl Repl, index int) (Repl, error){},
		FullSet:  set,
		Values:   map[string]lexml.Data{},
		Prompt:   "",
		Keywords: Keywords,
	}

	for {
		fmt.Printf("$_ ")
		repl.Prompt, err = reader.ReadString('\n')
		if err != nil {
			panic(fmt.Errorf("%s", err))
		}

		commands, err := repl.ParsePrompt()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("%+v", commands)
	}
}

func TestLaunch(set *lexml.Set) error {

	repl := Repl{
		Commands: map[string]func(repl Repl, index int) (Repl, error){},
		FullSet:  set,
		Values:   map[string]lexml.Data{},
		Prompt:   "",
		Keywords: Keywords,
	}

	var prompts = []string{
		"push",
		  "\"test\"",
		  "|",
		  "abcd",
		  "push push",
		  "push\"test\"",
		  "push|",
		  "push abcd",
		  "\"test\"push",
		"\"test\"\"test\"",
		"\"test\"|",
		"\"test\"abcd",
		"|push",
		"|\"test\"",
		"||",
		"|abcd",
	}

	// why the fuck this test doesn't work ??? It is acting like tokenizer.Repl.prompt is empty or something like that
	for _, ppt := range prompts {
		lexer := NewLexer([]byte(ppt))
		tokenizer := NewTokenizer(repl, &lexer)
		tokenizer.Tokenize()

		fmt.Printf("%s, %+v\n", ppt, tokenizer.Tokens)
	}
	return nil
}

func NewRepl() Repl {
	return Repl{
		Commands: map[string]func(repl Repl, index int) (Repl, error){},
		FullSet:  lexml.NewSet([]byte{}),
		Values:   make(map[string]lexml.Data),
		Prompt:   "",
		Keywords: Keywords,
	}
}
