package parser

import (
	"bufio"
	"fmt"
)

func MatchToken(tok int, buf *bufio.Reader) (err error) {
	if tok != C_tok.Tok {
		err = fmt.Errorf("Incorrect token %s, expected %c", C_tok.Value, tok)
		return
	}

	err = NextToken(buf)
	return
}

func NextToken(buf *bufio.Reader) (err error) {
	C_tok, err = getToken(buf)
	return
}

var peek = ' '

func getToken(buf *bufio.Reader) (tok Token, err error) {
WS:
	for {
		switch peek {
		case ' ', '\t', '\r', '\n':
		default:
			break WS
		}

		peek, _, err = buf.ReadRune()
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
			c, _, err = buf.ReadRune()
			if err != nil {
				return
			}

			if c == '\'' {
				break
			}

			tok.buf[bufLength] = c
			bufLength++
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		tok.Tok = INT
		for '0' <= peek && peek <= '9' {
			tok.buf[bufLength] = peek
			bufLength++

			peek, _, err = buf.ReadRune()
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
