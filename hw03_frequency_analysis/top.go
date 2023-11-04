package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const top = 10

func Top10(s string) []string {
	splited := strings.Fields(s)
	count := map[string]int{}
	for _, word := range splited {
		count[word]++
	}
	type pair struct {
		Key   string
		Value int
	}
	pairs := make([]pair, 0, len(count))
	for key, val := range count {
		pairs = append(pairs, pair{key, val})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Value != pairs[j].Value {
			return pairs[i].Value > pairs[j].Value
		}
		return pairs[i].Key < pairs[j].Key
	})
	result := make([]string, 0, top)
	if len(pairs) >= top {
		for _, val := range pairs[:top] {
			result = append(result, val.Key)
		}
	} else {
		for _, val := range pairs {
			result = append(result, val.Key)
		}
	}
	return result
}
