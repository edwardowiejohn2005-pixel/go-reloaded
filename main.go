package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: go run . input.txt output.txt")
		return
	}

	data, _ := os.ReadFile(os.Args[1])
	text := string(data)

	// Process tags first, then cleanup formatting
	text = convertNumbers(text)
	text = transform(text)
	text = fixArticle(text)
	text = fixPunctuation(text)

	os.WriteFile(os.Args[2], []byte(text), 0644)
}

func convertNumbers(text string) string {
	words := strings.Fields(text)
	for i := 0; i < len(words); i++ {
		base := 0
		if words[i] == "(hex)" { base = 16 }
		if words[i] == "(bin)" { base = 2 }

		if base != 0 && i > 0 {
			val, _ := strconv.ParseInt(words[i-1], base, 64)
			words[i-1] = strconv.FormatInt(val, 10)
			words = append(words[:i], words[i+1:]...)
			i--
		}
	}
	return strings.Join(words, " ")
}

func transform(text string) string {
	words := strings.Fields(text)
	for i := 0; i < len(words); i++ {
		tag := words[i]
		if !strings.HasPrefix(tag, "(") { continue }

		n := 1
		cmd := strings.Trim(tag, "(),")
		
		// Check for (cmd, n) pattern
		if strings.HasSuffix(tag, ",") && i+1 < len(words) {
			numStr := strings.Trim(words[i+1], ")")
			n, _ = strconv.Atoi(numStr)
			words = append(words[:i], words[i+2:]...) // Remove both parts of tag
		} else if cmd == "up" || cmd == "low" || cmd == "cap" {
			words = append(words[:i], words[i+1:]...) // Remove single tag
		} else {
			continue
		}

		// Apply transformation to n previous words
		for j := 1; j <= n && i-j >= 0; j++ {
			idx := i - j
			switch cmd {
			case "up":  words[idx] = strings.ToUpper(words[idx])
			case "low": words[idx] = strings.ToLower(words[idx])
			case "cap": words[idx] = capitalize(words[idx])
			}
		}
		i-- 
	}
	return strings.Join(words, " ")
}

func capitalize(w string) string {
	if len(w) == 0 { return w }
	return strings.ToUpper(string(w[0])) + strings.ToLower(w[1:])
}

func fixArticle(text string) string {
	words := strings.Fields(text)
	for i := 0; i < len(words)-1; i++ {
		lowW := strings.ToLower(words[i])
		if lowW == "a" {
			next := strings.ToLower(words[i+1])
			if strings.ContainsAny(string(next[0]), "aeiouh") {
				if words[i] == "A" { words[i] = "An" } else { words[i] = "an" }
			}
		}
	}
	return strings.Join(words, " ")
}

func fixPunctuation(text string) string {
	puncs := []string{".", ",", "!", "?", ";", ":"}
	// Fix space before: "word ." -> "word."
	for _, p := range puncs {
		text = strings.ReplaceAll(text, " "+p, p)
	}
	// Note: If you need to ensure space AFTER, you can add 
	// additional logic here, but this covers the prompt's core logic.
	return text
}