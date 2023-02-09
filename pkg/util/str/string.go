package str

import "strings"

func Repeat(char string, count int) string {
	builder := &strings.Builder{}
	for i := 1; i <= count; i++ {
		builder.WriteString(char)
	}
	return builder.String()
}
