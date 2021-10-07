/*
Generate a single json object that is a list of LAT patterns in the form:
[
{"LAT": <string>, "Pat": [<string>+]}
]

Though this works rather efficently, the concept behind it, though potentially
useful for populating a records based database, does not fit the needs of
the docuscope tagger as the resulting file is large and still needs processing
after loading into the tagger.
*/
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	
	"github.com/urfave/cli/v2"
	"gitlab.com/CMU_Sidecar/docuscope-dictionary-tools/docuscope-rules/internal/pkg/unobfuscate"
	"gitlab.com/CMU_Sidecar/docuscope-dictionary-tools/docuscope-rules/internal/pkg/fix"
	"gitlab.com/CMU_Sidecar/docuscope-dictionary-tools/docuscope-rules/internal/pkg/wordclasses"

	"golang.org/x/text/message"
)

func main() {
	var flagStats bool
	var cpuprofile string
	var memprofile string
	
	app := &cli.App{
		Name: "DocuScope Rule Database Generator",
		Usage: "Generates the JSON database from a directory containing LAT files and a _wordclasses.txt file.",
		UsageText: "docuscope-rules-db Dictionaries/default | gzip > default_db.json.gz",
		Version: "v1.0.0",
		Authors: []*cli.Author{
			&cli.Author{
				Name: "Michael Ringenberg",
				Email: unobfuscate.Unobfuscate("ringenbergATcmuDOTedu"),
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "stats",
				Usage: "Output statistics",
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
			return getDictionaryDB(c.Args().First(), flagStats)
			/*if _, err := os.Stdout.WriteString("["); err != nil {
				panic(err)
			}
			cnt, words, orig, missing, err := Lats(c.Args().First())
			if err == nil {
				w, errm := json.Marshal(Words{words})
				if errm != nil {
					panic(err)
				}
				if _, erro := os.Stdout.Write(w); erro != nil {
					panic(err)
				}
				if _, erro := os.Stdout.WriteString("]"); erro != nil {
					panic(err)
				}
				if flagStats {
					p := message.NewPrinter(message.MatchLanguage("en"))
					p.Fprintf(os.Stderr, "Rule Count: %d\n", cnt)
					p.Fprintf(os.Stderr, "Missing words added: %d; Original: %d; Final: %d\n", missing, orig, len(words))
				}
			}
			return err*/
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

type Rule struct {
	LAT string
	Pat []string
}
type Words struct {
	Words map[string][]string `json:"words"`
}

func getDictionaryDB(directory string, flagStats bool) error {
	words := make(map[string][]string)
	defaultWordsCount := 0
	missingWordsCount := 0
	ruleCount := 0
	wordclasses.ReadWords(words, filepath.Join(directory, "_wordclasses.txt"))
	defaultWordsCount = len(words)
	if _, err := os.Stdout.WriteString("["); err != nil {
		panic(err)
	}
	patternRe := regexp.MustCompile(`[!?\w'-]+|[!"#$%&'()*+,-./:;<=>?@[\]^_\` + "`" + `{|}~]`)
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error: unable to access %q: %v\n", path, err)
			panic(err)
		}
		base := filepath.Base(path)
		if !info.IsDir() && filepath.Ext(path) == ".txt" &&
			!strings.HasPrefix(base, "_") {
			lat := strings.TrimSuffix(base, ".txt")
			content, err := os.Open(filepath.Clean(path))
			if err != nil {
				fmt.Printf("Error: unable to access %q: %v\n", path, err)
				panic(err)
			}
			scanner := bufio.NewScanner(content)
			for scanner.Scan() {
				pattern := fix.Case(patternRe.FindAllString(scanner.Text(), -1))
				if (len(pattern) > 0) {
					b, err := json.Marshal(Rule{lat, pattern})
					if err != nil {
						panic(err)
					}
					ruleCount++
					if _, err := os.Stdout.Write(b); err != nil {
						panic(err)
					}
					if _, err := os.Stdout.WriteString(",\n"); err != nil {
						panic(err)
					}
					for _, w := range pattern {
						if wds, ok := words[w]; !ok {
							words[w] = append(wds, w)
							missingWordsCount++
						}
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
	w, err := json.Marshal(Words{words})
	if _, err := os.Stdout.Write(w); err != nil {
		panic(err)
	}
	if _, err := os.Stdout.WriteString("]"); err != nil {
		panic(err)
	}
	if flagStats {
		p := message.NewPrinter(message.MatchLanguage("en"))
		p.Fprintf(os.Stderr, "Rule Count: %d\n", ruleCount)
		p.Fprintf(os.Stderr, "Missing words added: %d; Original: %d; Final: %d\n", missingWordsCount, defaultWordsCount, len(words))
	}
	return nil
}

/**
 * asyncronusly walk all of the directories
 */
func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)
	go func() {
		defer close(paths)
		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: unable to access %q: %v\n", path, err)
				return err
			}
			if !info.Mode().IsRegular() || info.IsDir() ||
				filepath.Ext(path) != ".txt" ||
				strings.HasPrefix(filepath.Base(path), "_") {
				// fmt.Fprintln(os.Stderr, "Ignoring: ", path)
				return nil
			}
			select {
			case paths <- path:
			case <-done:
				return errors.New("File directory walk canceled")
			}
			return nil
		})
	}()
	return paths, errc
}

type result struct {
	path string
	lat Rule
	err error
}

/**
 * Extract rules from a LAT file.
 */
func ruliser(done <-chan struct{}, paths <-chan string, c chan<- result) {
	patternRe := regexp.MustCompile(`[!?\w'-]+|[!"#$%&'()*+,-./:;<=>?@[\]^_\` + "`" + `{|}~]`)
	for path := range paths {
		base := filepath.Base(path)
		content, err := os.Open(filepath.Clean(path))
		if err != nil {
			select {
			case c <- result{path, Rule{}, err}:
			case <-done:
				return
			}
		}
		scanner := bufio.NewScanner(content)
		for scanner.Scan() {
			pattern := fix.Case(patternRe.FindAllString(scanner.Text(), -1))
			if (len(pattern) > 0) {
				lat := strings.TrimSuffix(base, ".txt")
				select {
				case c <- result{path, Rule{lat, pattern}, err}:
				case <-done:
					return
				}
			}
		}
		if err := content.Close(); err != nil {
			log.Fatal("Could not close content file: ", err)
			c <- result{path, Rule{}, err}
		}
	}
}

func Lats(root string) (int, map[string][]string, int, int, error) {
	words := make(map[string][]string)
	defaultWordsCount := 0
	missingWordsCount := 0
	wordclasses.ReadWords(words, filepath.Join(root, "_wordclasses.txt"))
	defaultWordsCount = len(words)
	
	done := make(chan struct{})
	defer close(done)
	paths, errc := walkFiles(done, root)

	c := make(chan result)
	var wg sync.WaitGroup
	const numDigesters = 30
	wg.Add(numDigesters)
	for i := 0; i < numDigesters; i++ {
		go func() {
			ruliser(done, paths, c)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(c)
	}()

	ruleCount := 0
	for r := range c {
		if r.err != nil {
			return ruleCount, words, defaultWordsCount, missingWordsCount, r.err
		}
		// Add any missing words
		for _, w := range r.lat.Pat {
			if wds, ok := words[w]; !ok {
				words[w] = append(wds, w)
				missingWordsCount++
			}
		}
		b, err := json.Marshal(r.lat)
		if err != nil {
			return ruleCount, words, defaultWordsCount, missingWordsCount, err
		}
		ruleCount += 1
		if _, err := os.Stdout.Write(b); err != nil {
			return ruleCount, words, defaultWordsCount, missingWordsCount, err
		}
		if _, err := os.Stdout.WriteString(",\n"); err != nil {
			return ruleCount, words, defaultWordsCount, missingWordsCount, err
		}
	}
	if err := <-errc; err != nil {
		return ruleCount, words, defaultWordsCount, missingWordsCount, err
	}
	return ruleCount, words, defaultWordsCount, missingWordsCount, nil
}
