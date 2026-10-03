package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"sync"
)

func collect() error {
	type sourceResult struct {
		source string
		links  []string
		err    error
	}
	results := make([]sourceResult, len(sources))
	var wg sync.WaitGroup
	for i, source := range sources {
		wg.Add(1)
		go func(i int, source string) {
			defer wg.Done()
			links, err := fetchSource(source)
			results[i] = sourceResult{source, links, err}
		}(i, source)
	}
	wg.Wait()

	seen := map[string]bool{}
	var candidates []string
	for _, result := range results {
		if result.err != nil {
			fmt.Printf("skip %s: %v\n", result.source, result.err)
			continue
		}
		rand.Shuffle(len(result.links), func(i, j int) {
			result.links[i], result.links[j] = result.links[j], result.links[i]
		})
		added := 0
		for _, link := range result.links {
			if added >= MaxPerSource {
				break
			}
			if key := connectionKey(link); !seen[key] {
				seen[key] = true
				candidates = append(candidates, link)
				added++
			}
		}
		fmt.Printf("%4d candidates from %s\n", added, result.source)
	}
	if len(candidates) == 0 {
		return errors.New("no candidates collected from any source")
	}
	return os.WriteFile(CandidatesFile, []byte(strings.Join(candidates, "\n")+"\n"), 0644)
}
