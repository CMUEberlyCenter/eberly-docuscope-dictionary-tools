package wordclasses

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/**
 * Reads _wordclasses.txt file associated with a DocuScope dictionary.
 *
 * @param words: the map of word class to array of members.
 * @param wordclassesPath: location of the _wordclasses.txt file.
 */
func ReadWords(words map[string][]string, wordclassesPath string) {
	curClass := "NONE"
	wordclasses, err := os.Open(filepath.Clean(wordclassesPath))
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(wordclasses)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		switch len(line) {
		case 1:
			word := strings.ToLower(line[0])
			_, ok := words[word]
			if !ok {
				words[word] = append(words[word], word)
			}
			words[word] = pushnew(words[word], curClass)
		case 2:
			curClass = "!" + strings.ToUpper(line[1])
		default:
			// noop
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if err := wordclasses.Close(); err != nil {
		log.Fatal("Could not close word classes file: ", err)
	}
}

// append only if not already an element of the slice.
func pushnew(slice []string, val string) []string {
	for _, ele := range slice {
		if ele == val {
			return slice
		}
	}
	return append(slice, val)
}
