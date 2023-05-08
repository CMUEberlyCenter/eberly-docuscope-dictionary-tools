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
	      "word2": {
	        "Cat": [[],["word3"]]
	      }
	    }
	  },
	  "shortRules": {
	    "word": "Cat"
	  },
	  "words": {
	    "!BANG": ["bang"]
	  }
	}
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
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/urfave/cli/v2"

	"gitlab.com/CMU_Sidecar/docuscope-dictionary-tools/docuscope-rules/internal/pkg/fix"
	"gitlab.com/CMU_Sidecar/docuscope-dictionary-tools/docuscope-rules/internal/pkg/unobfuscate"
	"gitlab.com/CMU_Sidecar/docuscope-dictionary-tools/docuscope-rules/internal/pkg/wordclasses"
)

type RulesMap map[string]map[string]map[string][][]string

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
	Rules      RulesMap            `json:"rules"`
	ShortRules map[string]string   `json:"shortRules"`
	Words      map[string][]string `json:"words"`
}

/**
 * Add a pattern to the rules map.
 */
func add(m RulesMap, lat string, rule []string) {
	mm, ok := m[rule[0]]
	if !ok {
		mm = make(map[string]map[string][][]string)
		m[rule[0]] = mm
	}
	mmm, ok := mm[rule[1]]
	if !ok {
		mmm = make(map[string][][]string)
		mm[rule[1]] = mmm
	}
	mmm[lat] = append(mmm[lat], rule[2:])
}

func genDictionaryRules(directory string, flagStats bool) error {
	rules := make(RulesMap)
	shortRules := make(map[string]string)
	words := make(map[string][]string)
	missingWordsCount := 0
	defaultWordsCount := 0
	patternRe := regexp.MustCompile(`[!?\w'-]+|[!"#$%&'()*+,-./:;<=>?@[\]^_\` + "`" + `{|}~]`)

	wordclasses.ReadWords(words, filepath.Join(directory, "_wordclasses.txt"))
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
			content, err := os.Open(filepath.Clean(path))
			if err != nil {
				panic(err)
			}

			scanner := bufio.NewScanner(content)
			for scanner.Scan() {
				pattern := fix.Case(patternRe.FindAllString(scanner.Text(), -1))
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
			if err := content.Close(); err != nil {
				log.Fatal("Could not close content file: ", err)
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
	if _, err := os.Stdout.Write(b); err != nil {
		panic(err)
	}
	return nil
}

func main() {
	var flagStats bool
	var cpuprofile string
	var memprofile string

	app := &cli.App{
		Name:      "DocuScope Rule File Generator",
		Usage:     "Generates the JSON rules file from a directory containing LAT files and a _wordclasses.txt file.",
		UsageText: "docuscope-rules Dictionaries/default | gzip > default.json.gz",
		Version:   "v1.0.4",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Michael Ringenberg",
				Email: unobfuscate.Unobfuscate("ringenbergATcmuDOTedu"),
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "stats",
				Usage:       "Output statistics",
				Destination: &flagStats,
			},
			&cli.StringFlag{
				Name:        "cpuprofile",
				Value:       "",
				Usage:       "Write cpu profile to `file`",
				Destination: &cpuprofile,
			},
			&cli.StringFlag{
				Name:        "memprofile",
				Value:       "",
				Usage:       "Write memory profile to `file`",
				Destination: &memprofile,
			},
		},
		Action: func(c *cli.Context) error {
			return genDictionaryRules(c.Args().First(), flagStats)
		},
	}

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal("Could not create CPU profile: ", err)
		}
		//defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Could not start CPU profile: ", err)
		}
		//defer pprof.StopCPUProfile()
		if rerr := app.Run(os.Args); rerr != nil {
			log.Fatal(rerr)
		}

		pprof.StopCPUProfile()
		if err := f.Close(); err != nil {
			log.Fatal("Could not close cpu profile: ", err)
		}
	} else {
		err := app.Run(os.Args)
		if err != nil {
			log.Fatal(err)
		}
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Fatal("Could not create memory profile: ", err)
		}
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("Could not write memory profile: ", err)
		}
		if err := f.Close(); err != nil {
			log.Fatal("Could not close memory profile: ", err)
		}
	}
}
