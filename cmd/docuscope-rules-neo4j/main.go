/*
default:
real	300m2.788s
user	9m13.238s
sys	7m44.774s
*/
package main

import(
	"bufio"
	//"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/golobby/dotenv"
	"github.com/urfave/cli/v2"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"gitlab.com/CMU_Sidecar/docuscope-dictionary-tools/docuscope-rules/internal/pkg/unobfuscate"
	"gitlab.com/CMU_Sidecar/docuscope-dictionary-tools/docuscope-rules/internal/pkg/fix"
	"gitlab.com/CMU_Sidecar/docuscope-dictionary-tools/docuscope-rules/internal/pkg/wordclasses"
)

type Env struct {
	Neo4J struct {
		Uri string `env:"NEO4J_URI"`
		User string `env:"NEO4J_USER"`
		Pass string `env:"NEO4J_PASSWORD"`
	}
}

func main() {
	var flagStats bool
	var cpuprofile string
	var memprofile string

	config := Env{}
	file, err := os.Open(".env")
	if err != nil {
		log.Fatal("Could not open .env: ", err)
	}
	
	err = dotenv.NewDecoder(file).Decode(&config);
	if err != nil {
		log.Fatal("Could not decode .env: ", err)
	}

	app := &cli.App{
		Name: "DocuScope Rules for Neo4j",
		Usage: "Populates a neo4j database with the rule and wordclasses.",
		UsageText: "docuscope-rule-neo4j Dictionaries/default",
		Version: "v0.0.1",
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
			return addDictionary(c.Args().First(),
				config.Neo4J.Uri, config.Neo4J.User,
				config.Neo4J.Pass, flagStats)
		},
	}
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal("Could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Could not start CPU profile: ", err)
		}
		if rerr := app.Run(os.Args); rerr != nil {
			log.Fatal(rerr)
		}

		pprof.StopCPUProfile()
		if err := f.Close(); err != nil {
			log.Fatal("Could not close cpu profile: ", err)
		}
	} else {
		if err := app.Run(os.Args); err != nil {
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

type MemoizedQuery func(int) string
func memoQuery() MemoizedQuery {
	cache := make(map[int]string)
	cache[0] = ""
	return func(index int) string {
		if val, found := cache[index]; found {
			return val
		}
		//fmt.Printf("Generating query for %d\n", index)
		var qry strings.Builder
		qry.WriteString("MERGE (s0:Start {word: $p0}) ")
		for j := 1; j < index; j++ {
			fmt.Fprintf(&qry, "MERGE (s%d)-[:NEXT {word: $p%d}]->(s%d) ", j-1, j, j)
		}
		qry.WriteString("MERGE (l:Lat {lat: $lat}) ")
		fmt.Fprintf(&qry, "MERGE (s%d)-[:LAT]->(l);", index-1)
		result := qry.String()
		cache[index] = result
		return result
	}
}

func addDictionary(directory string, uri string, username string, password string, flagStats bool) error {
	fmt.Printf("Connecting to %q as %q.\n", uri, username) 
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatal("Could not open database: ", uri, username, err)
		panic(err)
	}
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	// Create index
	_, txerr := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run("CREATE INDEX start_index IF NOT EXISTS FOR (s:Start) ON (s.word);", map[string]interface{}{})
		if err != nil {
			fmt.Printf("Error on start_index: %v.\n", err)
			return nil, err
		}
		_, err1 := tx.Run("CREATE INDEX lat_index IF NOT EXISTS FOR (l:Lat) ON (l.lat);", map[string]interface{}{})
		if err1 != nil {
			fmt.Printf("Error on lat_index: %v.\n", err1)
			return nil, err1
		}
		_, err2 := tx.Run("CREATE INDEX next_index IF NOT EXISTS FOR ()-[n:NEXT]->() ON (n.word);", map[string]interface{}{})
		if err2 != nil {
			fmt.Printf("Error on next_index: %v\n", err2)
		}
		return nil, err2
	})
	if txerr != nil {
		fmt.Printf("Error on index transaction: %v\n", txerr)
		panic(txerr)
	}
	// Start memoized query provider.
	merges := memoQuery()
	/*for i := 1; i<len(merges); i++ {
		merges[i] = qry.String()
	}*/

	words := make(map[string][]string)
	defaultWordsCount := 0
	missingWordsCount := 0
	//ruleCount := 0
	wordclasses.ReadWords(words, filepath.Join(directory, "_wordclasses.txt"))
	defaultWordsCount = len(words)
	patternRe := regexp.MustCompile(`[!?\w'-]+|[!"#$%&'()*+,-./:;<=>?@[\]^_\` + "`" + `{|}~]`)
	walkerr := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error: unable to access %q: %v\n", path, err)
			panic(err)
		}
		base := filepath.Base(path)
		if !info.IsDir() && filepath.Ext(path) == ".txt" &&
			!strings.HasPrefix(base, "_") {
			lat := strings.TrimSuffix(base, ".txt")
			content, oerr := os.Open(filepath.Clean(path))
			if oerr != nil {
				fmt.Printf("Error: unable to access %q: %v\n", path, oerr)
				panic(oerr)
			}
			numPatterns := 0
			scanner := bufio.NewScanner(content)
			_, txerr := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
				for scanner.Scan() {
					pattern := fix.Case(patternRe.FindAllString(scanner.Text(), -1))
					if (len(pattern) > 0) {
						var pmap = map[string]interface{} {
							"lat": lat,
						}
						for i, v := range(pattern) {
							pmap[fmt.Sprint("p", i)] = v
						}
						numPatterns++
						if numPatterns % 1000 == 0 {
							fmt.Printf("\r%d", numPatterns)
						}
						_, err := transaction.Run(
							merges(len(pattern)),
							pmap)
						// add length binning count
						if err != nil {
							fmt.Printf("Query error: %q %d %v\n", lat, len(pattern), pattern, err)
							fmt.Printf("Query: %q %v\n", merges(len(pattern)), pmap)
							return nil, err
						}
					}
					for _, w := range pattern {
						if wds, ok := words[w]; !ok {
							words[w] = append(wds, w)
							missingWordsCount++
						}
					}
				}
				return nil, nil
			})			
			if txerr != nil {
				fmt.Printf("Error on transaction: %q: %v\n", lat, txerr)
				panic(txerr)
			}
			if err := content.Close(); err != nil {
				log.Fatal("Could not close content file: ", err)
				panic(err)
			}
			fmt.Printf("\r%q %d\n", lat, numPatterns) //, time
		}
		return nil
	})
	if walkerr != nil {
		panic(walkerr)
	}
	//w, err := json.Marshal(words)
	//if err != nil {
	//	panic(err)
	//}
	//if _, err := os.Stdout.Write(w); err != nil {
	//	panic(err)
	//}
	if flagStats {
		fmt.Fprintln(os.Stderr, "Missing words:", defaultWordsCount,
			missingWordsCount, len(words))
	}
	return nil
}
