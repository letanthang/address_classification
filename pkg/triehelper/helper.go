package triehelper

import (
	"address_classification/entity"
	"address_classification/pkg/stringutil"
	"address_classification/trie"
	"address_classification/trie/parse"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strings"
)

func ClassifyAddress(input string, trieDic *trie.Trie, reversedTrie *trie.Trie) entity.Result {
	input = NormalizeInput(input)
	result := parse.DynamicParseWithSkipV2(input, trieDic, reversedTrie)
	return result
}

func NormalizeInput(input string) string {
	input = strings.ToLower(input)

	input = stringutil.RemoveVietnameseAccents(input)

	return input
}

func ImportWardDB(filename string) []entity.Ward {
	results := []entity.Ward{}

	// Mở file CSV
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Tạo một CSV reader
	reader := csv.NewReader(file)
	reader.Comma = ';'

	// Skip csv header
	_, err = reader.Read()
	if err != nil {
		log.Fatal("Unable to read header:", err)
	}

	// Duyệt qua từng dòng và in ra
	i := 1
	for {
		record, err := reader.Read()
		if err != nil {
			// Nếu không còn dòng nào hoặc gặp lỗi
			break
		}

		ward := entity.Ward{
			Province:     record[0],
			ProvinceCode: record[1],
			District:     record[2],
			DistrictCode: record[3],
			Name:         record[4],
			Code:         record[5],
		}

		i++
		results = append(results, ward)
	}

	return results
}

func ImportTestCases(fileName string) []entity.TestCase {
	var testCases []entity.TestCase

	bytes, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(bytes, &testCases); err != nil {
		log.Fatal(err)
	}

	return testCases

}
