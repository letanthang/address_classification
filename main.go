package main

import (
	"address_classification/entity"
	"address_classification/pkg/triehelper"
	"address_classification/trie"
	"log"
)

func main() {
	wards := triehelper.ImportWardDB("./assets/wards.csv")
	trieDic := trie.NewTrie()
	trieDic.BuildTrieWithWards(wards)

	//testCases := importTestCases("./assets/inputs.json")

	input := []string{
		//"nguyen tri phuong, phuong 10, quan 10, tp ho chi minh",
		//"nguyen tri, phuong 10, quan 1, tp ho chi minh",
		//"nguyen tri, phuong 100, quan 11, tp ho chi minh",
		//"nguyen tri phuong 100 quan 111 tp ho chi minh",
		//"quan 111 tp ho chi minh",
		//"tp ho chí minh quận 2", // missing Q2 in db/trie
		"Xa Quảng Thọ,T.P Sầm Swn,TY. thanh Hóa",
	}

	//word := trieDic.ExtractWord(input[0], 17)
	//log.Println(word)
	for _, address := range input {
		result := triehelper.ClassifyAddress(address, trieDic)
		logResult(result)
	}
}

func logResult(result entity.Result) {
	log.Printf("Province: %s, District: %s, Ward: %s\n", result.Province, result.District, result.Ward)
}
