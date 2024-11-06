package main

import (
	"address_classification/entity"
	"address_classification/pkg/triehelper"
	"address_classification/trie"
	"address_classification/trie/parse"
	"fmt"
	"log"
)

func main() {
	TestWithRealCases()
}

func TestSimple() {
	parse.Debug = true

	wards := triehelper.ImportWardDB("./assets/wards.csv")
	trieDic := trie.NewTrie(false)
	trieDic.BuildTrieWithWards(wards)

	reversedTrie := trie.NewTrie(true)
	reversedTrie.BuildTrieWithWards(wards)

	input := []string{

		"46/8F Trung Chánh 2 Trung Chánh, Hóc Môn, TP. Hồ Chí Minh",
	}

	for _, address := range input {
		result := triehelper.ClassifyAddress(address, trieDic, reversedTrie)
		logResult(result)
	}
}

func TestWithRealCases() {
	parse.Debug = true
	wards := triehelper.ImportWardDB("./assets/wards.csv")

	trieTree := trie.NewTrie(false)
	trieTree.BuildTrieWithWards(wards)

	reversedTrie := trie.NewTrie(true)
	reversedTrie.BuildTrieWithWards(wards)

	cases := triehelper.ImportTestCases("./assets/inputs.json")

	for i, c := range cases {
		result := triehelper.ClassifyAddress(c.Input, trieTree, reversedTrie)
		if result.Ward != c.Output.Ward || result.District != c.Output.District || result.Province != c.Output.Province {
			//fmt.Println(i, "Failed")
			logResult(result)
		} else {
			fmt.Println(i, "Passed")
		}
	}
}

func logResult(result entity.Result) {
	log.Printf("Result : Province %s, District %s, Ward %s\n", result.Province, result.District, result.Ward)
}
