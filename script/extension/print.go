package extension

import (
	"fmt"
)

// Only for egvm script debugging
func Println(a ...any) {
	fmt.Println(a...)
}

func Printf(format string, a ...any) {
	fmt.Printf(format, a...)
}
