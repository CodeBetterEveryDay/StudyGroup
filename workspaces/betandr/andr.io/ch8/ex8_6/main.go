// Add depth-limiting to the concurrent crawler. That is, if the user sets
// `-depth=3`, then only URLs reachable by at most three links will be fetched.
// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

var depth = flag.Int("depth", 3, "Only URLs reachable by this number of links will be fetched")

type link struct {
	URL   string
	Depth int
}

func crawl(url link) []string {
	fmt.Println(url)
	list, err := links.Extract(url.URL)
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	worklist := make(chan []string) // lists of URLs, may have duplicates
	unseenLinks := make(chan link)  // de-duplicated URLs

	// Add command-line arguments to worklist.
	go func() { worklist <- os.Args[1:] }()

	// Create 20 crawler goroutines to fetch each unseen link.
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				// increment depth here
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	seen := make(map[string]bool)
	for list := range worklist {
		for _, l := range list {
			if !seen[l] {
				seen[l] = true
				unseenLinks <- link{l, 0}
			}
		}
	}
}
