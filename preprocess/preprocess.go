package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

type Pokemon struct {
	No   int    `json:"no"`
	Name string `json:"name"`
}

type Data struct {
	Data []Pokemon `json:"data"`
}

type Word struct {
	Pokemon
	Ruby   string `json:"reading"`
	Start  string `json:"start"`
	End    string `json:"end"`
	IsLast bool   `json:"is_last"`
}

type SmallKana string

const (
	SmallYa  SmallKana = "ャ"
	SmallYu  SmallKana = "ュ"
	SmallYo  SmallKana = "ョ"
	SmallA   SmallKana = "ァ"
	SmallI   SmallKana = "ィ"
	SmallU   SmallKana = "ゥ"
	SmallE   SmallKana = "ェ"
	SmallO   SmallKana = "ォ"
	SmallTsu SmallKana = "ッ"
)

var largeKanaMap = map[SmallKana]string{}

func init() {
	largeKanaMap[SmallYa] = "ヤ"
	largeKanaMap[SmallYu] = "ユ"
	largeKanaMap[SmallYo] = "ヨ"
	largeKanaMap[SmallA] = "ア"
	largeKanaMap[SmallI] = "イ"
	largeKanaMap[SmallU] = "ウ"
	largeKanaMap[SmallE] = "エ"
	largeKanaMap[SmallO] = "オ"
	largeKanaMap[SmallTsu] = "ツ"
}

func main() {
	f, err := os.Open("data/zukan.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var data Data
	json.NewDecoder(f).Decode(&data)

	words := []Word{}
	for _, p := range data.Data {
		var w Word
		w.Pokemon = p

		var ruby string
		switch p.Name {
		case "ニドラン♂":
			ruby = "ニドラオス"
		case "ニドラン♀":
			ruby = "ニドランメス"
		case "ポリゴン2":
			ruby = "ポリゴンツー"
		case "ポリゴンZ":
			ruby = "ポリゴンゼット"
		default:
			ruby = p.Name
		}
		w.Ruby = ruby

		r := []rune(w.Ruby)
		w.Start = string(r[:1])
		endIndex := len(r) - 1
		end := r[endIndex : endIndex+1]

		pattern := `[アイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワンヴガギグゲゴザジズゼゾダヂヅデドバビブベボパピプペポ]`
		re, err := regexp.Compile(pattern)
		if err != nil {
			panic(err)
		}

		// ルールは(https://pokemon-shiritori.com/howto)に準拠
		// 伸ばし棒で終わる場合
		if string(end) == "ー" {
			endIndex--
			end = r[endIndex : endIndex+1]
		}

		// 小文字で終わる場合
		if l, ok := largeKanaMap[SmallKana(end)]; ok {
			end = []rune(l)
		}

		w.End = string(end)

		if len(w.End) == 0 {
			fmt.Println(p.Name)
			panic("end is empty")
		}

		if w.End == "ン" {
			w.IsLast = true
		}

		if !re.MatchString(w.End) {
			fmt.Println(p.Name)
			panic("不正な文字が含まれています")
		}

		words = append(words, w)
	}

	j, err := json.Marshal(words)
	if err != nil {
		panic(err)
	}

	f, err = os.Create("data/words.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(j)
	if err != nil {
		panic(err)
	}
}
