package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

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
	var pairs []pair
	for key, val := range count {
		pairs = append(pairs, pair{key, val})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Value != pairs[j].Value {
			return pairs[i].Value > pairs[j].Value
		}
		return pairs[i].Key < pairs[j].Key
	})
	var result []string
	for _, val := range pairs {
		result = append(result, val.Key)
	}
	if len(result) >= 10 {
		return result[:10]
	}
	return result
}
