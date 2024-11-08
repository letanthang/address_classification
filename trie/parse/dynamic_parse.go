package parse

import (
	"address_classification/entity"
	"address_classification/trie"
	"fmt"
	"slices"
	"sort"
	"strings"
)

var (
	Debug           bool
	DebugFlag       string
	Words           []string
	OriginLocations []entity.Location
	SkipWords       []string
	CorrectedResult entity.Result
	Locations       []entity.Location
	delimiters      = []rune{' ', ',', '-'}
)

func DynamicParse(originSentence string, trieDic *trie.Trie, reversedTrie *trie.Trie) entity.Result {
	var (
		skipWords []string
		words     []string
		locations []entity.Location
	)

	DebugFlag = "empty"
	result := entity.Result{}

	if trieDic == nil || len(originSentence) == 0 {
		return result
	}

	wordMap := map[string]struct{}{}

	var extract func(sentence string)
	extract = func(sentence string) {
		if len(sentence) == 0 {
			return
		}
		first := 0
		offset := 0
		i := 0
		// skip delimiters
		for {
			if !slices.Contains(delimiters, rune(sentence[first])) {
				break
			}
			first++
		}
		sentence = sentence[first:]

		// skip word not in trie
		skip := trieDic.Skip(sentence)

		if skip >= len(sentence) {
			return
		}

		if skip > 0 {
			skipWords = append(skipWords, sentence[:skip])
		}

		sentence = sentence[skip:]

		//DebugFlag = DebugFlag + " " + sentence
		for offset < len(sentence) {
			i++
			word, node := trieDic.ExtractWord(sentence, offset)
			offset = offset + len(word)

			if word == "" {
				return
			}

			if _, ok := wordMap[word]; !ok {
				wordMap[word] = struct{}{}
				words = append(words, word)
				locations = append(locations, node.Locations...)
			}

			offset = len(word)

			if offset >= len(sentence) {
				return
			}

			extract(sentence[offset:])
		}

		return
	}

	extract(originSentence)
	printWords(words, "words")
	printWords(skipWords, "skips")

	Words = words
	OriginLocations = locations
	locations = trie.FilterLocation(locations, words, originSentence)
	Locations = locations
	SkipWords = skipWords

	result = getLocationFromLocations(locations)
	if result.IsComplete() {
		return result
	}

	CorrectedResult = DynamicParseWithLevenshtein(skipWords, reversedTrie)
	mergeResult(&CorrectedResult, &result)

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

			locations := make([]entity.Location, len(node.Locations))
			copy(locations, node.Locations)
			sort.Sort(entity.Locations(locations))
			AddLocationToResult(&result, node.Locations[0])
		}

		printWords(correctedWords, "corrected words: ")
	}
	return result

}

func getLocationFromLocations(locations []entity.Location) entity.Result {
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

func mergeResult(source, destination *entity.Result) {
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
