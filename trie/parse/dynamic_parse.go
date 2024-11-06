package parse

import (
	"address_classification/entity"
	"address_classification/trie"
	"fmt"
	"slices"
	"sort"
	"strings"
)

var Debug bool

var delimiters = []rune{' ', ',', '-'}

var dynamicStatus = map[string]string{}

func getDynamicStatus(sentence string, offset int) string {
	return dynamicStatus[sentence]
}

func setDynamicStatus(sentence string, offset int, word string) {
	dynamicStatus[sentence] = word
}

func DynamicParseWithSkipV2(originSentence string, trieDic *trie.Trie, reversedTrie *trie.Trie) entity.Result {
	result := entity.Result{}
	skipWords := []string{}
	locations := []entity.Location{}
	firstAttempt := true

	if trieDic == nil || len(originSentence) == 0 {
		return result
	}

	var extract func(sentence string) (bool, []string)
	extract = func(sentence string) (bool, []string) {
		if len(sentence) == 0 {
			return true, nil
		}
		first := 0
		offset := 0
		// skip delimiters
		for {
			if !slices.Contains(delimiters, rune(sentence[first])) {
				break
			}
			first++
		}

		sentence = sentence[first:]

		for offset < len(sentence) {
			word, node, skip := trieDic.ExtractWordWithSkipping(sentence, offset)
			if skip > 0 && firstAttempt {
				skipWords = append(skipWords, sentence[0:skip-1])
			}

			if word == "" {
				firstAttempt = false
				return false, nil
			}
			locations = append(locations, node.Locations...)

			offset = skip + len(word)

			if offset >= len(sentence) {
				firstAttempt = false
				return true, []string{word}
			}

			ok, words := extract(sentence[offset:])

			if ok {
				words = append(words, word)
				return true, words
			}
		}

		return false, nil
	}

	_, words := extract(originSentence)
	printWords(words, "words")
	printWords(skipWords, "skips")
	locations = trie.FilterLocation(locations, words)
	//fmt.Println(entity.Locations(locations).ToString())

	result = GetLocationFromLocations(locations)
	if result.IsComplete() {
		return result
	}

	correctedResult := DynamicParseWithLevenshtein(skipWords, reversedTrie)
	MergeResult(&correctedResult, &result)

	return result
}

func DynamicParseWithLevenshtein(skipWords []string, trieDic *trie.Trie) entity.Result {
	result := entity.Result{}
	if len(skipWords) == 0 || trieDic == nil {
		return result
	}

	correctedWords := []string{}
	for _, skipWord := range skipWords {
		correctedWord, _, node := trieDic.ExtractWordWithAutoCorrect(skipWord)
		if correctedWord != "" {
			correctedWords = append(correctedWords, correctedWord)

			sort.Sort(entity.Locations(node.Locations))
			AddLocationToResult(&result, node.Locations[0])
		}

		printWords(correctedWords, "corrected words: ")
	}
	return result

}

func GetLocationFromLocations(locations []entity.Location) entity.Result {
	result := entity.Result{}

	for _, location := range locations {
		AddLocationToResult(&result, location)
	}
	return result
}

func AddLocationToResult(r *entity.Result, location entity.Location) {
	switch location.LocationType {
	case entity.LocationTypeWard:
		r.Ward = trie.WardMap[location.ID].NoPrefixName
	case entity.LocationTypeDistrict:
		r.District = trie.DistrictMap[location.ID].NoPrefixName
	case entity.LocationTypeProvince:
		r.Province = trie.ProvinceMap[location.ID].NoPrefixName
	}
}

func MergeResult(source, destination *entity.Result) {
	if source.Ward != "" && destination.Ward == "" {
		destination.Ward = source.Ward
	}

	if source.District != "" && destination.District == "" {
		destination.District = source.District
	}

	if source.Province != "" && destination.Province == "" {
		destination.Province = source.Province
	}
}

func printWords(words []string, wordType string) {
	if !Debug {
		return
	}

	text := strings.Join(words, "|")
	fmt.Println(wordType + ": " + text)
}
