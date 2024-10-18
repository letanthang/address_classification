package parse

import (
	"address_classification/entity"
	"address_classification/trie"
	"fmt"
	"slices"
	"strings"
)

var delimiters = []rune{' ', ',', '-'}

var dynamicStatus = map[string]string{}

func getDynamicStatus(sentence string, offset int) string {
	return dynamicStatus[sentence]
}

func setDynamicStatus(sentence string, offset int, word string) {
	dynamicStatus[sentence] = word
}

func DynamicParseWithSkipV2(originSentence string, trieDic *trie.Trie) entity.Result {
	result := entity.Result{}
	locations := []entity.Location{}
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
			locations = append(locations, node.Locations...)

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

	_, words := extract(first, end)
	printWords(words)
	locations = trie.FilterLocation(locations, words)
	//fmt.Println(entity.Locations(locations).ToString())

	result = GetLocationFromLocations(locations)
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

func printWords(words []string) {
	text := strings.Join(words, "|")
	fmt.Println("Words: " + text)
}
