package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"unicode/utf8"
)

//
// Replace any fancy prompt with just a '$' character
//
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

//
// Remove any right-side prompt (currently timestamp and timing info)
//
func stripTimestamp(buf []byte) ([]byte, error) {
	s := string(buf)

	re := regexp.MustCompile(`(  )+([^ ]+)?( . )?\d\d:\d\d:\d\d$`)
	s = re.ReplaceAllString(s, "")

	return []byte(s), nil
}

//
// If an application name was supplied on the command line, we can add a custom
// prefix and suffix based on the application.  I use this to wrap the paste in
// Markdown syntax when pasting into Slack or Discord
//
func wrappers(app string) (prefix, suffix string) {
	switch app {
	case "Slack":
	case "Discord":
		prefix, suffix = "```\n", "```"
	}
	return
}

func main() {
	fgPtr := flag.String("app", "", "Foreground App Name")
	flag.Parse()

	foreground := *fgPtr

	prefix, suffix := wrappers(foreground)

	fmt.Printf(prefix)

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

	fmt.Printf(suffix)
}
