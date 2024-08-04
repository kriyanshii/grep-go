package main

import (
	// Uncomment this to pass the first stage
	// "bytes"

	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(string(line), pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

// func matchLine(line string, pattern string) (bool, error) {
// 	if utf8.RuneCountInString(pattern) == 0 {
// 		return false, fmt.Errorf("unsupported pattern: %q", pattern)
// 	}

// 	var ok bool

// 	// You can use print statements as follows for debugging, they'll be visible when running tests.
// 	fmt.Println("Logs from your program will appear here!")

// 	// Uncomment this to pass the first stage
// 	if pattern == "\\d" {
// 		ok = bytes.ContainsAny(line, "1234567890")
// 	} else if pattern == "\\w" {
// 		ok = bytes.ContainsAny(line, "abcdefghijklmnopqrstuvwxyz")
// 	} else if pattern[0] == '[' && pattern[len(pattern)-1] == ']' {
// 		n := len(pattern)
// 		if pattern[1] == '^' {
// 			ok = !bytes.ContainsAny(line, pattern[2:n-1])
// 		} else {
// 			ok = bytes.ContainsAny(line, pattern[1:n-1])
// 		}
// 	} else {
// 		ok = bytes.ContainsAny(line, pattern)
// 	}

// 	return ok, nil
// }

func matchLine(line string, pattern string) (bool, error) {
	if pattern[0] == '^' {
		return matchPattern(line, pattern[1:], 0), nil
	}
	for i := 0; i < len(line); i++ {
		if matchPattern(line, pattern, i) {
			return true, nil
		}
	}
	return false, nil
}

func matchPattern(line string, pattern string, pos int) bool {
	n, m := len(pattern), len(line)
	j := pos
	for i := 0; i < n; i++ {
		if j >= m {
			return pattern[i] == '$'
		}
		if pattern[i] == '\\' && i+1 < n {
			if pattern[i+1] == 'd' && !unicode.IsDigit(rune(line[j])) {
				return false
			} else if pattern[i+1] == 'w' && !unicode.IsLetter(rune(line[j])) {
				return false
			} else {
				i++
			}
		} else if pattern[i] == '[' && i+1 < n && pattern[i+1] == '^' {
			endPos := strings.Index(pattern[i:], "]")
			matchThisPattern := pattern[i:endPos]
			if strings.Contains(matchThisPattern, string(line[j])) {
				return false
			}
			i = endPos
		} else if pattern[i] == '[' {
			endPos := strings.Index(pattern[i:], "]")
			matchThisPattern := pattern[i:endPos]
			if !strings.Contains(matchThisPattern, string(line[j])) {
				return false
			}
			i = endPos
		} else if pattern[i] == '+' {
			if i == 0 {
				fmt.Println("Invalid pattern")
				return false
			}
			for j < m && line[j] == pattern[i-1] {
				j++
			}
			j--
		} else {
			if j < m && line[j] != pattern[i] {
				return false
			}
		}
		j++
	}
	return true
}
