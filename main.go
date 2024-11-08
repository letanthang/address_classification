package main

import (
	"address_classification/entity"
	"address_classification/pkg/triehelper"
	"address_classification/trie"
	"address_classification/trie/parse"
	"fmt"
	"log"
	"strings"
)

func main() {
	testMode := 0

	switch testMode {
	case 1:
		TestSimple()
	case 2:
		TestWithRealCases()
	default:
		DebugTrie()
	}
}

func DebugTrie() {
	wards := triehelper.ImportWardDB("./assets/wards.csv")
	trieDic := trie.NewTrie(false)
	trieDic.BuildTrieWithWards(wards)
	_, node := trieDic.ExtractWord("pho", 0)
	fmt.Println(node)
	fmt.Println(trie.WardMap[node.Locations[0].ID])
}

func TestSimple() {
	//parse.Debug = true

	wards := triehelper.ImportWardDB("./assets/wards.csv")
	trieDic := trie.NewTrie(false)
	trieDic.BuildTrieWithWards(wards)

	reversedTrie := trie.NewTrie(true)
	reversedTrie.BuildTrieWithWards(wards)

	input := []string{
		"Khu phố 4 Thị trấn, Dương Minh Châu, Tây Ninh",
	}

	for i := 0; i < 100000; i++ {
		address := input[0]
		result := triehelper.ClassifyAddress(address, trieDic, reversedTrie)
		if i == 0 {
			//logResult(result)
		}

		if result.Ward == "Hồ" {
			logResult(result)
			logResult(parse.CorrectedResult)
			fmt.Println("skip words", parse.SkipWords)
			fmt.Println(parse.DebugFlag)
			fmt.Println("words", parse.Words)
			fmt.Println(entity.Locations(parse.OriginLocations).ToString())
			fmt.Println(entity.Locations(parse.Locations).ToString())
			break
		}
	}
}

func TestWithRealCases() {
	//parse.Debug = true
	wards := triehelper.ImportWardDB("./assets/wards.csv")

	trieTree := trie.NewTrie(false)
	trieTree.BuildTrieWithWards(wards)

	reversedTrie := trie.NewTrie(true)
	reversedTrie.BuildTrieWithWards(wards)

	cases := triehelper.ImportTestCases("./assets/inputs.json")

	for i, c := range cases {
		result := triehelper.ClassifyAddress(c.Input, trieTree, reversedTrie)
		if result.Ward != c.Output.Ward || result.District != c.Output.District || result.Province != c.Output.Province {
			logResult(result)
			fmt.Println(parse.DebugFlag)
			logWords(parse.Words)
			fmt.Println(entity.Locations(parse.OriginLocations).ToString())
			fmt.Println(entity.Locations(parse.Locations).ToString())
			return
		} else {
			fmt.Println(i, "Passed")
		}
	}
}

func logResult(result entity.Result) {
	log.Printf("Result : Province %s, District %s, Ward %s\n", result.Province, result.District, result.Ward)
}

func logWords(words []string) {
	text := strings.Join(words, "|")
	fmt.Println("words: ", text)
}
