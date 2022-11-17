// Classifier attempts to classify documents by topic.
package cirrus

import (
	"sort"

	"github.com/hvlck/txt"
)

type Topic struct {
	Name   string
	Weight float64
}

type pair struct {
	word       string
	occurences int
}

type pairList []pair

func (p pairList) Len() int           { return len(p) }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pairList) Less(i, j int) bool { return p[i].occurences < p[j].occurences }

func Classify(text string) []Topic {
	_, err := Recognize(text)
	if err != nil {
		panic(err)
	}

	words := make(map[string]int)

	tokens := txt.Tokenize(text, txt.DefaultSplitter, nil, txt.TokenizerStopwords)
	for _, v := range tokens {
		if occurences, ok := words[v]; ok {
			words[v] = occurences + 1
		} else {
			words[v] = 1
		}
	}

	pairs := make(pairList, 0, len(words))

	for k, v := range words {
		pairs = append(pairs, pair{word: k, occurences: v})
	}

	sort.Sort(pairs)

	return make([]Topic, 0)
}
