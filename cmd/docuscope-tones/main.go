/*
Tool for generating DocuScope json tones files.

Usage: docuscope_tones < _tones.txt > tones.json

JSON Schema: see api/docuscope_tones_schema.json
*/
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"gitlab.com/CMU_Sidecar/docuscope-dictionary-tools/docuscope-rules/internal/pkg/unobfuscate"
)

func add(m map[string]map[string][]string, cluster string, dimension string, lats []string) {
	mm, ok := m[cluster]
	if !ok {
		mm = make(map[string][]string)
		m[cluster] = mm
	}
	// pushnew lats onto existing (no duplicates).
	// This handles the problem where a given tone is repeated.
	// This should probably be broader to check for no lat duplicates
	// as that will cause errors in docuscope-tag as it will complain
	// that indicies should be unique.
	for _, lat := range lats {
		to_add := true
		for _, ele := range mm[dimension] {
			if ele == lat {
				to_add = false
				break
			}
		}
		if to_add {
			mm[dimension] = append(mm[dimension], lat)
		}
	}
}

func main() {
	app := &cli.App{
		Name:      "DocuScope Tones Converter",
		Usage:     "Convert a DocuScope _tone.txt file to json.",
		UsageText: "docuscope-tones < _tones.txt > tones.json",
		Version:   "v1.0.0",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Michael Ringenberg",
				Email: unobfuscate.Unobfuscate("ringenbergATcmuDOTedu"),
			},
		},
		Action: func(c *cli.Context) error {
			return tonesToJson()
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func tonesToJson() error {
	var cluster string
	var dimension string
	tones := make(map[string]map[string][]string)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) > 1 {
			switch line[0] {
			case "CLUSTER:":
				cluster = line[1]
			case "DIMENSION:":
				dimension = line[1]
			case "LAT:", "LAT*:", "CLASS:":
				//add to tones
				add(tones, cluster, dimension, line[1:])
			default:
				//noop
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading standard input:", err)
	}
	b, err := json.Marshal(tones)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error outputting:", err)
	}
	if _, err := os.Stdout.Write(b); err != nil {
		log.Fatal(err)
	}
	return nil
}
