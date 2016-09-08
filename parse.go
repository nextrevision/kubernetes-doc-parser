package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	css "github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

var (
	path           = flag.String("path", "", "file path to parse (supports wildcards)")
	pattern        = flag.String("pattern", "", "css selector pattern")
	attr           = flag.String("attr", "href", "attribute to select")
	ignore         = flag.String("ignore", "", "regexp of matches to ignore")
	verbose        = flag.Bool("verbose", false, "print verbose output")
	replacePattern = flag.String("replace-pattern", "", "replacement regexp")
	replaceString  = flag.String("replace-string", "", "replacement string")
	results        = []string{}
)

func main() {
	flag.Parse()

	compiledIgnore := regexp.MustCompile(*ignore)
	compiledReplacePattern := regexp.MustCompile(*replacePattern)

	if *path == "" {
		log.Fatal("Must specify a path with -path")
	}

	if *pattern == "" {
		log.Fatal("Must specify a pattern with -pattern")
	}

	if *ignore != "" {

	}

	paths, err := filepath.Glob(*path)
	if err != nil {
		log.Fatalf("Error parsing path '%s': %v", *path, err)
	}

	compiledPattern, err := css.Compile(*pattern)

	for _, p := range paths {
		if *verbose {
			log.Printf("Parsing %s", p)
		}

		f, err := os.Open(p)
		if err != nil {
			log.Fatalf("Error opening path '%s': %v", p, err)
		}

		parsedHTML, err := html.Parse(f)
		if err != nil {
			log.Fatalf("Error parsing file '%s' as HTML: %v", p, err)
		}

		for _, match := range compiledPattern.MatchAll(parsedHTML) {
			for _, a := range match.Attr {
				if a.Key == *attr {
					val := strings.TrimSuffix(a.Val, "/")

					if *ignore != "" && compiledIgnore.MatchString(val) {
						continue
					}

					if *replacePattern != "" {
						val = compiledReplacePattern.ReplaceAllString(val, *replaceString)
					}

					if *verbose {
						log.Printf("Found match: %s", val)
					}

					if !strings.HasSuffix(val, ".html") {
						if !exists(fmt.Sprintf("%s/index.html")) {
							log.Printf("Found something odd, no index.html found for %s", val)
						}
						val = fmt.Sprintf("%s/index.html", val)
					}

					results = append(results, val)
				}
			}
		}
	}
	payload, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("Error converting results to JSON: %v", err)
	}

	println(string(payload))
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
