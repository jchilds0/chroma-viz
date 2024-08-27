package parser

import (
	"log"
	"strconv"
	"strings"
)

func AddAttribute(b strings.Builder, name string, value any) {
	b.WriteRune('\'')
	b.WriteString(name)
	b.WriteString("': ")

	switch v := value.(type) {
	case int:
		b.WriteString(strconv.Itoa(v))
	case int64:
		b.WriteString(strconv.FormatInt(v, 10))
	case float64:
		b.WriteString(strconv.FormatFloat(v, 'f', 64, 10))
	case string:
		b.WriteRune('\'')
		b.WriteString(v)
		b.WriteRune('\'')
	default:
		log.Printf("Unknown attribute type %v", value)
	}
}
