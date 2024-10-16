package parse

import (
	"address_classification/entity"
	"address_classification/trie"
	"slices"
)

var delimiters = []rune{' ', ',', '-'}

var dynamicStatus = map[string]string{}

func getDynamicStatus(sentence string, offset int) string {
	return dynamicStatus[sentence]
}

func setDynamicStatus(sentence string, offset int, word string) {
	dynamicStatus[sentence] = word
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
			word, _ := trieDic.ExtractWord(sentence, offset)
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

func DynamicParseWithSkip(originSentence string, trieDic *trie.Trie) (bool, []string) {
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
			word, _, skip := trieDic.ExtractWordWithSkipping(sentence, offset)
			if word == "" {
				return false, nil
			}

			offset = len(word) + skip
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

func DynamicParseWithSkipV2(originSentence string, trieDic *trie.Trie) entity.Result {
	result := entity.Result{}
	first := 0
	end := len(originSentence) - 1

	if trieDic == nil || end <= 0 {
		return result
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
			word, node, skip := trieDic.ExtractWordWithSkipping(sentence, offset)
			if word == "" {
				return false, nil
			}

			AddNodeToResult(&result, node)

			offset = len(word) + skip
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

	extract(first, end)

	return result
}

func AddNodeToResult(r *entity.Result, node *trie.Node) {

	if node == nil || !node.IsEnd || r.IsComplete() {
		return
	}

	switch node.Type {
	case trie.NodeTypeWard:
		r.Ward = trie.WardMap[node.IDs[0]].Name
	case trie.NodeTypeDistrict:
		r.District = trie.DistrictMap[node.IDs[0]].Name
	case trie.NodeTypeProvince:
		r.Province = trie.ProvinceMap[node.IDs[0]].Name
	}
}
