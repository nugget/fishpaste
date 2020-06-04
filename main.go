package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"unicode/utf8"
)

func replacePrompt(buf []byte) ([]byte, error) {
	newPrompt := []byte("$ ")

	promptEnd := 0
	promptPossible := 0

	b := buf
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		promptPossible += size
		// fmt.Printf("%c (%d) %v [%d]\n", r, r, size, promptPossible)
		if r == 57520 {
			promptEnd = promptPossible + 1
		}
		b = b[size:]
	}

	// fmt.Println("promptEnd", promptEnd, "len", len(buf))

	if promptEnd > 0 && promptEnd < len(buf) {
		buf = buf[promptEnd:]
		buf = append(newPrompt, buf...)
	} else if promptEnd > len(buf) {
		buf = newPrompt
	}

	return buf, nil
}

func stripTimestamp(buf []byte) ([]byte, error) {
	s := string(buf)

	re := regexp.MustCompile(`(  )+([^ ]+)?( . )?\d\d:\d\d:\d\d$`)
	s = re.ReplaceAllString(s, "")

	return []byte(s), nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var err error

		buf := scanner.Bytes()

		buf, err = replacePrompt(buf)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		buf, err = stripTimestamp(buf)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(string(buf))
	}
}
