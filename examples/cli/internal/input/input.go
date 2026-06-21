package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Confirm() bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		s, _ := reader.ReadString('\n')
		s = strings.TrimSuffix(s, "\n")
		s = strings.ToLower(s)

		if s == "n" {
			return false
		} else if s == "y" {
			break
		} else {
			_, _ = fmt.Fprintln(os.Stderr, "Please enter Y or N")
			continue
		}
	}

	return true
}
