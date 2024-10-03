package parse

import (
	"address_classification/trie"
	"slices"
)

var delimiters = []rune{' ', ',', '-'}

var dynamicStatus = map[int]bool{}

func makeKey(first, end int) int {
	return first<<16 + end
}

func DynamicParse(originSentence string, trieDic *trie.Trie) (bool, []string) {
	first := 0
	end := len(originSentence) - 1

	if trieDic == nil || end <= 0 {
		return false, nil
	}

	var extract func(first, end int) (bool, []string)
	extract = func(first, end int) (bool, []string) {
		if first >= end {
			return true, nil
		}
		offset := 0
		// skip delimiters
		for {
			if !slices.Contains(delimiters, rune(originSentence[first])) {
				break
			}
			first++
		}

		sentence := originSentence[first:]

		for offset < len(sentence) {
			word := trieDic.ExtractWord(sentence, offset)
			if word == "" {
				return false, nil
			}

			offset = len(word)
			nextFirst := first + offset

			if nextFirst == len(originSentence)-1 {
				return true, []string{word}
			}

			ok, words := extract(nextFirst, end)

			if ok {
				words = append(words, word)
				return true, words
			}
		}

		return false, nil
	}

	return extract(first, end)
}
