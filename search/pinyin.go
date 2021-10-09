package search

import (
	"strings"

	"github.com/go-ego/gpy"
)

func PinYin(text string) string {
	tags := strings.Split(text, ",")
	rTags := make([]string, 0, 2*len(tags))

	for _, tag := range tags {
		words := gpy.LazyConvert(tag, nil)
		var firsts = make([]string, 0, len(words))

		for _, word := range words {
			firsts = append(firsts, strings.Split(word, "")[0])
		}

		if len(words) == 0 {
			continue
		}
		rTags = append(rTags, strings.Join(words, ""), strings.Join(firsts, ""))
	}

	return strings.Join(rTags, ",")
}
