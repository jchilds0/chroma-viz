package parser

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	s := "-10 'hello' '123' 0 {'x': 12345};"
	val := []string{"-10", "hello", "123", "0", "", "x", "", "12345", ""}
	tok := []int{INT, STRING, STRING, INT, '{', STRING, ':', INT, '}'}

	buf := strings.NewReader(s)
	err := NextToken(buf)
	if err != nil {
		t.Error(err)
	}

	for i := range val {
		if C_tok.Value != val[i] {
			t.Fatal("Incorrect token ", C_tok.Value, " expected ", val[i])
		}

		err = MatchToken(tok[i], buf)
		if err != nil {
			t.Fatal(err)
		}
	}
}
