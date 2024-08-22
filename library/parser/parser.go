package parser

import (
	"fmt"
	"io"
	"log"
)

// tokens
const (
	INT = iota + 256
	STRING
)

type Token struct {
	Tok   int
	Value string
	buf   []rune
}

var C_tok Token
var b []rune = make([]rune, 0, 100)

func MatchToken(tok int, buf io.RuneReader) (err error) {
	if tok != C_tok.Tok {
		log.Printf("Buffer: %s", string(b))
		err = fmt.Errorf("Incorrect token %s, expected %c", C_tok.Value, tok)
		log.Println(err)
		return
	}

	err = NextToken(buf)
	return
}

func NextToken(buf io.RuneReader) (err error) {
	C_tok, err = getToken(buf)
	return
}

var peek = ' '

func getToken(buf io.RuneReader) (tok Token, err error) {
WS:
	for {
		switch peek {
		case ' ', '\t', '\r', '\n':
		default:
			break WS
		}

		peek, _, err = readRune(buf)
		if err != nil {
			return
		}
	}

	tok.buf = make([]rune, 64)
	bufLength := 0

	switch peek {
	case '\'':
		var c rune
		tok.Tok = STRING
		peek = ' '

		for {
			c, _, err = readRune(buf)
			if err != nil {
				return
			}

			if c == '\'' {
				break
			}

			tok.buf[bufLength] = c
			bufLength++
		}
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		tok.Tok = INT

		if peek == '-' {
			tok.buf[bufLength] = peek
			bufLength++

			peek, _, err = readRune(buf)
			if err != nil {
				return
			}
		}

		for '0' <= peek && peek <= '9' {
			tok.buf[bufLength] = peek
			bufLength++

			peek, _, err = readRune(buf)
			if err != nil {
				return
			}
		}
	default:
		tok.Tok = int(peek)
		tok.buf[0] = peek
		peek = ' '
	}

	tok.Value = string(tok.buf[:bufLength])
	return
}

func readRune(buf io.RuneReader) (r rune, n int, err error) {
	r, n, err = buf.ReadRune()

	b = append(b, r)
	if len(b) == 100 {
		b = b[1:]
	}

	return
}
