package mineralrights

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readIntWithPrompt(prompt string, v Validator) int64 {
	fmt.Printf("%s ", prompt)

	i := readInt()

	// keep reading until the validation passes
	for {
		if v(i) {
			return i
		}

		fmt.Printf("? ")
		i = readInt()
	}

	return i
}

func readInt() int64 {
	for {
		input := readString()

		v := strings.TrimSpace(input)

		if len(v) == 0 {
			return 0
		}

		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}

		return i
	}
}

func readString() string {
	reader := bufio.NewReader(os.Stdin)

	// ignoring errors reading from keyboard input
	s, _ := reader.ReadString('\n')

	return s
}
