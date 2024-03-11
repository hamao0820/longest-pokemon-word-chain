//go:build ignore

package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

//go:embed data/words.json
var wordsJSON []byte

type Pokemon struct {
	No   int    `json:"no"`
	Name string `json:"name"`
}

type Word struct {
	Pokemon
	Ruby   string `json:"reading"`
	Start  string `json:"start"`
	End    string `json:"end"`
	IsLast bool   `json:"is_last"`
}

var words []Word
var noMap map[int]Word = make(map[int]Word)
var startMap map[string][]Word = make(map[string][]Word)

func init() {
	json.Unmarshal(wordsJSON, &words)

	for _, w := range words {
		noMap[w.No] = w
		startMap[w.Start] = append(startMap[w.Start], w)
	}
}

type Edge struct {
	From int
	To   int
}

func main() {
	li := map[int][]Word{} // 隣接リスト
	for i := 0; i < len(words); i++ {
		w := words[i]
		if !isFistGeneration(&w.Pokemon) {
			continue
		}
		for _, v := range startMap[w.End] {
			if !isFistGeneration(&v.Pokemon) {
				continue
			}
			li[w.No] = append(li[w.No], v)
		}
	}

	getLongestPath := func(start int) ([]int, uint64) {
		visited := make(map[int]bool)
		path := []int{}
		longestPath := []int{}
		maxLength := 0
		var count uint64
		// 最長単純路を求める
		var backtrack func(int) bool
		backtrack = func(v int) bool {
			count++
			visited[v] = true
			path = append(path, v)
			if noMap[v].IsLast && len(path) > maxLength {
				maxLength = len(path)
				longestPath = append([]int{}, path...)
			}
			for _, u := range li[v] {
				if visited[u.No] {
					continue
				}
				backtrack(u.No)
			}
			path = path[:len(path)-1]
			visited[v] = false
			return false
		}
		backtrack(start)
		return longestPath, count
	}

	type result struct {
		Start  Word    `json:"start"`
		Path   []Word `json:"chain_words"`
		Length int    `json:"length"`
	}
	results := []result{}
	f := func(i int, wg *sync.WaitGroup) {
		defer wg.Done()
		s := time.Now()
		p, _ := getLongestPath(i)
		chainWords := []Word{}
		for _, v := range p {
			chainWords = append(chainWords, noMap[v])
		}
		results = append(results, result{noMap[i], chainWords, len(p)})
		fmt.Printf("start: %d; time: %v\n", i, time.Since(s))
	}
	var wg sync.WaitGroup
	for i := 1; i <= 151; i++ {
		wg.Add(1)
		go f(i, &wg)
	}
	wg.Wait()
	fmt.Println("done")

	// results to JSON
	resultsJSON, _ := json.Marshal(results)
	resultFile, _ := os.Create("data/results.json")
	defer resultFile.Close()
	resultFile.Write(resultsJSON)
}

func isFistGeneration(p *Pokemon) bool {
	return p.No <= 151
}
