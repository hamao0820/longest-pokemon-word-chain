package main

import (
	_ "embed"
	"encoding/json"
	"syscall/js"
)

//go:embed data/results.json
var resultsJSON []byte

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

type Result struct {
	Start  Word   `json:"start"`
	Path   []Word `json:"chain_words"`
	Length int    `json:"length"`
}

var results []Result
var words []Word
var noMap map[int]Word = make(map[int]Word)

func init() {
	json.Unmarshal(resultsJSON, &results)
	json.Unmarshal(wordsJSON, &words)
	for _, w := range words {
		noMap[w.No] = w
	}
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("goGetLongestChain", js.FuncOf(func(_ js.Value, args []js.Value) any {
		start := args[0].Int()
		var res Result
		for _, r := range results {
			if r.Start.No == start {
				res = r
				break
			}
		}
		arr := js.Global().Get("Array").New(res.Length)
		for i, w := range res.Path {
			arr.SetIndex(i, js.ValueOf(w.No))
		}
		return arr
	}))

	<-c
}
