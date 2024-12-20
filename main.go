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
	testMode := 2

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
	//fmt.Println(trieDic.IsEnd("trang dai"))
	//_, node := trieDic.ExtractWord("trang dai", 0)
	//fmt.Println(node)
	//fmt.Println(trie.WardMap[node.Locations[0].ID])

	sentence := "thị trấn Ba Hàng"

	sentence = triehelper.NormalizeInput(sentence)

	word, _ := trieDic.ExtractWord(sentence, 0)
	fmt.Println(word)
}

func TestSimple() {
	//parse.Debug = true

	wards := triehelper.ImportWardDB("./assets/wards.csv")
	trieDic := trie.NewTrie(false)
	trieDic.BuildTrieWithWards(wards)

	reversedTrie := trie.NewTrie(true)
	reversedTrie.BuildTrieWithWards(wards)

	input := []string{
		"Tiểu khu 3, thị trấn Ba Hàng, huyện Phổ Yên, tỉnh Thái Nguyên.",
	}

	for i := 0; i < 1; i++ {
		address := input[0]
		result := triehelper.ClassifyAddress(address, trieDic, reversedTrie)
		//if i == 0 {
		//	logResult(result)
		//}

		if true {
			logResult(result)
			printResult(parse.CorrectedResult)
			logWords(parse.Words)
			logWords(parse.SkipWords)
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

	failNum := 0
	for i, c := range cases {
		result := triehelper.ClassifyAddress(c.Input, trieTree, reversedTrie)
		if result.Ward != c.Output.Ward || result.District != c.Output.District || result.Province != c.Output.Province {
			//logResult(result)
			//fmt.Println(parse.DebugFlag)
			//logWords(parse.Words)
			//logWords(parse.SkipWords)
			//fmt.Println(entity.Locations(parse.OriginLocations).ToString())
			//fmt.Println(entity.Locations(parse.Locations).ToString())
			//return
			failNum++
		} else {
			fmt.Println(i, "Passed")
		}
	}

	fmt.Println("Fail num: ", failNum, "/", len(cases))
}

func logResult(result entity.Result) {
	log.Printf("Result : Province %s, District %s, Ward %s\n", result.Province, result.District, result.Ward)
}

func printResult(result entity.Result) {
	fmt.Printf("Result : Province %s, District %s, Ward %s\n", result.Province, result.District, result.Ward)
}

func logWords(words []string) {
	text := strings.Join(words, "|")
	fmt.Println("words: ", text)
}
