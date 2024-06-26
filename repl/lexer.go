package repl

import (
    "io"
)

type Lexer struct {
    Offset int
    Bytes []byte
}

func (lexer *Lexer) Next() bool {

    lexer.Offset = lexer.Offset + 1
    if lexer.Offset == len(lexer.Bytes) {
        return false
    }
    return true
}

func (lexer *Lexer) Peek() (byte, error) {
    if lexer.Offset <= len(lexer.Bytes) {
        return 0, io.EOF
    }
    return lexer.Bytes[lexer.Offset + 1], nil
}

func (lexer *Lexer) Byte() byte {
    return lexer.Bytes[lexer.Offset]
}

func (lexer *Lexer) Text() string {
    return string(lexer.Byte())
}

func NewLexer(bytes []byte) Lexer {
    return Lexer{
        Offset: 0,
        Bytes: bytes,
    }
}
