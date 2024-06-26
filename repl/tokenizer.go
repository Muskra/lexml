package repl

import (
// "fmt"
)

const (
	NONE uint8 = iota
	PIPE
	STRING
	KEYWORD
	COMMAND
	LABEL
)

type Tokenizer struct {
	Repl     Repl
	Scanner  *Lexer
	Tokens   []Token
	Offset   int
	HasBound bool
}

type Token struct {
	Kind  uint8
	Value string
}

var Keywords = []string{
	"pop",   // delete a value from the Repl.Values map
	"exit",  // terminates the program
	"help",  // display help function
	"print", // prints a value
	"push",  // push a value to the Repl.Values map
}

// Tokenize function returns a list of tokens from a string.
// this function is not yet complete, it seems that an empty prefix and an empty suffix are added to the slice. I thought the fix was to simply slice the slice from [1:] or [:len(ofSlice)-1] but this makes the program crash for some reason.
func (tokenizer *Tokenizer) Tokenize() {

	var check bool = true

	for check {

		if len(tokenizer.Tokens) == 0 {
			tokenizer.Tokens = append(tokenizer.Tokens, NewToken(NONE))
		}
        // not sure if it's mandatory
		/*if tokenizer.Scanner.Offset >= len(tokenizer.Scanner.Bytes) {
			break
		}*/

		switch tokenizer.Scanner.Text() {
		case "\"":
			lookupString(tokenizer)
			tokenizer.Offset = tokenizer.Offset + 1
		case "|":
			tokenizer.Tokens = append(
				tokenizer.Tokens,
				Token{
					Kind:  PIPE,
					Value: "|",
				},
			)
			tokenizer.Offset = tokenizer.Offset + 1
		default:
			tokenizer.Tokens[tokenizer.Offset].Value = tokenizer.Tokens[tokenizer.Offset].Value + tokenizer.Scanner.Text()
		}

		check = tokenizer.Scanner.Next()
	}

	tokenizer.Trim()

	//tokenizer.CheckSpecialTokens()
}

// this seems overcomplicated for the problem that it is
func (tokenizer *Tokenizer) Trim() {

	switch tokenizer.Tokens[0].Value {
	case "":
		tokenizer.Tokens = tokenizer.Tokens[1:]
	case " ":
		tokenizer.Tokens = tokenizer.Tokens[1:]
	}

	switch tokenizer.Tokens[len(tokenizer.Tokens)-1].Value {
	case "":
		tokenizer.Tokens = tokenizer.Tokens[:len(tokenizer.Tokens)-1]
	case " ":
		tokenizer.Tokens = tokenizer.Tokens[:len(tokenizer.Tokens)-1]
	}
}

/*
// here we could check for any COMMAND or KEYWORD and then if nothing match, we assign a LABEL Kind to it.
func (tokenizer *Tokenizer) CheckSpecialTokens() {

    for index, token := range tokenizer.Tokens {
        if token.Kind == NONE {

        }
    }

}*/

func NewTokenizer(repl Repl, lexer *Lexer) Tokenizer {

	return Tokenizer{
		Repl:     repl,
		Scanner:  lexer,
		Tokens:   make([]Token, 0),
		Offset:   0,
		HasBound: false,
	}
}

func NewToken(kind uint8) Token {
	return Token{
		Kind:  kind,
		Value: "",
	}
}

func lookupString(tokenizer *Tokenizer) {

	tokenizer.Tokens = append(tokenizer.Tokens, NewToken(STRING))

	for tokenizer.Scanner.Text() != "\"" {
		tokenizer.Tokens[tokenizer.Offset].Value = tokenizer.Tokens[tokenizer.Offset].Value + tokenizer.Scanner.Text()
	}
}

// isCommand function checks if a value as string is a command. It's checking within a Repl.Commands map Type. Returns a boolean Type
func isCommand(repl Repl, value string) bool {
	commands := repl.Commands
	if _, ok := commands[value]; ok {
		return true
	} else {
		return false
	}
}

// isKeyword function checks if a value as string is a keyword from the Keywords public variable. Returns a boolean Type
func isKeyword(value string) bool {
	for _, keyw := range Keywords {
		if value == keyw {
			return true
		}
	}
	return false
}
