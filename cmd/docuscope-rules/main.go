/*
Generate a single JSON representation of all of the LAT files in a directory.

Usage:
> docuscope_rules Dictionaries/default > rules.json
> docuscope_rules Dictionaries/default | gzip > rules.json.gz

JSON schema: See api/docuscope_rules_schema.json

Example:
{
  "rules": {
    "word1": {
      "word2": [["Cat", ["word1", "word2"]]]
    }
  },
  "shortRules": {
    "word": "Cat"
  },
  "words": {
    "!BANG": ["bang"]
  }
}

TODO: change rules to {<word>:{<word>:{LAT: <lat>, patterns: [[<string-2>+]]}}}
*/
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"
)

/**
 * Corrects letter case for words and wordclasses.
 * Words should be lowercase.
 * Wordclasses, indicated by ! prefix, should be uppercase.
 */
func fixCase(pat []string) []string {
	ret := make([]string, len(pat))
	for i, v := range pat {
		if strings.HasPrefix(v, "!") {
			ret[i] = strings.ToUpper(v)
		} else {
			ret[i] = strings.ToLower(v)
		}
	}
	return ret
}

/*
DocuScopeDictionary contains the patterns and words in the dictionary used by
DocuScope to parse a document.
The rules are organized such that it is a partial reverse lookup with the
initial bigram serving as the indicies to an array of arrays where the first
element is the id of the LAT and the second is an array with the full pattern.
The ShortRules are a mapping of unigram to LAT id.
Words is a mapping of !CLASS or words to an array of words or classes.
*/
type DocuScopeDictionary struct {
	Rules      map[string]map[string][]interface{} `json:"rules"`
	ShortRules map[string]string                   `json:"shortRules"`
	Words      map[string][]string                 `json:"words"`
}

/**
 * Add a pattern to the rules map.
 */
func add(m map[string]map[string][]interface{}, lat string, rule []string) {
	mm, ok := m[rule[0]]
	if !ok {
		mm = make(map[string][]interface{})
		m[rule[0]] = mm
	}
	r := []interface{}{lat, rule}
	mm[rule[1]] = append(mm[rule[1]], r)
}

func genDictionaryRules(directory string, flagStats bool) error {
	rules := make(map[string]map[string][]interface{})
	shortRules := make(map[string]string)
	words := make(map[string][]string)
	missingWordsCount := 0
	defaultWordsCount := 0
	patternRe := regexp.MustCompile(`[!?\w'-]+|[!"#$%&'()*+,-./:;<=>?@[\]^_\` + "`" + `{|}~]`)

	readWords(words, filepath.Join(directory, "_wordclasses.txt"))
	defaultWordsCount = len(words)
	err := filepath.Walk(directory, func(path string,
		info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error: unable to access %q: %v\n",
				path, err)
			panic(err)
		}
		base := filepath.Base(path)
		if !info.IsDir() && filepath.Ext(path) == ".txt" &&
			!strings.HasPrefix(base, "_") {
			lat := strings.TrimSuffix(base, ".txt")
			content, err := os.Open(path)
			defer content.Close()
			if err != nil {
				panic(err)
			}

			scanner := bufio.NewScanner(content)
			for scanner.Scan() {
				pattern := fixCase(patternRe.FindAllString(scanner.Text(), -1))
				switch len(pattern) {
				case 0:
					//noop
				case 1:
					shortRules[pattern[0]] = lat
				default:
					add(rules, lat, pattern)
				}
				for _, w := range pattern {
					if wds, ok := words[w]; !ok {
						words[w] = append(wds, w)
						missingWordsCount++
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	if flagStats {
		fmt.Fprintln(os.Stderr, "Missing words:", defaultWordsCount,
			missingWordsCount, len(words))
	}

	b, err := json.Marshal(DocuScopeDictionary{rules, shortRules, words})
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(b)
	return nil
}

func readWords(words map[string][]string, wordclassesPath string) {
	curClass := "NONE"
	wordclasses, err := os.Open(wordclassesPath)
	defer wordclasses.Close()
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
			words[word] = appendIfMissing(words[word], curClass)
		case 2:
			curClass = "!" + strings.ToUpper(line[1])
		default:
			//noop
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

// append only if no already an element of the slice.
func appendIfMissing(slice []string, val string) []string {
	for _, ele := range slice {
		if ele == val {
			return slice
		}
	}
	return append(slice, val)
}

func unobfuscate(email string) string {
	reAT := regexp.MustCompile(`AT`)
	reDOT := regexp.MustCompile(`DOT`)
	out := reAT.ReplaceAllString(email, "@")
	out = reDOT.ReplaceAllString(out, ".")
	return out
}

func main() {
	var flagStats bool

	app := &cli.App{
		Name:      "DocuScope Rule File Generator",
		Usage:     "Generates the JSON rules file from a directory containing LAT files and a _wordclasses.txt file.",
		UsageText: "docuscope-rules Dictionaries/default | gzip > default.json.gz",
		Version:   "v1.0.3",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Michael Ringenberg",
				Email: unobfuscate("ringenbergATcmuDOTedu"),
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "stats",
				Usage:       "Output statistics",
				Destination: &flagStats,
			},
		},
		Action: func(c *cli.Context) error {
			return genDictionaryRules(c.Args().First(), flagStats)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
