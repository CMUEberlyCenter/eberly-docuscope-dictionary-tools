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
	"gitlab.com/CMU_Sidecar/docuscope-dictionary-tools/docuscope-rules/internal/pkg/unobfuscate"
	"gitlab.com/CMU_Sidecar/docuscope-dictionary-tools/docuscope-rules/internal/pkg/fix"
	"gitlab.com/CMU_Sidecar/docuscope-dictionary-tools/docuscope-rules/internal/pkg/wordclasses"
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
					if _, err := os.Stdout.Write(b); err != nil {
						panic(err)
					}
					if _, err := os.Stdout.WriteString(","); err != nil {
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
		fmt.Fprintln(os.Stderr, "Missing words:", defaultWordsCount, missingWordsCount, len(words))
	}
	return nil
}
