package main

import (
	"address_classification/entity"
	"address_classification/pkg/stringutil"
	"address_classification/trie"
	"address_classification/trie/parse"
	"bufio"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strings"
)

type TestCase struct {
	Input  string `json:"text"`
	Output Result `json:"result"`
}

type Result struct {
	Ward     string `json:"ward"`
	District string `json:"district"`
	Province string `json:"province"`
}

func main() {

	wards := importWardDB("./assets/wards.csv")

	trieDic := trie.NewTrie()
	trieDic.BuildTrieWithWards(wards)

	//trieDic.Print()

	//testCases := importTestCases("./assets/inputs.json")

	input := []string{
		"nguyen tri phuong, phuong 10, quan 10, tp ho chi minh",
		//"nguyen tri phuong, phuong 10, quan 10, tp ho chi minh",
		//"nguyen tri phuong, phuong 10, tp ho chi minh, quan 10",
		//"nguyen tri phuong phuong 10 tp ho chi minh quan 10",
	}

	//word := trieDic.ExtractWord(input[0], 17)
	//log.Println(word)

	for _, address := range input {
		result := classifyAddress(normalizeInput(address), trieDic)
		log.Println("final result", result)
	}

}

func classifyAddress(input string, trieDic *trie.Trie) Result {
	result := Result{}

	ok, words := parse.DynamicParseWithSkip(input, trieDic)
	if ok {
		logWords(words)
	}

	return result
}

func normalizeInput(input string) string {
	input = strings.ToLower(input)

	input = stringutil.RemoveVietnameseAccents(input)

	return input
}

func importWardDB(filename string) []entity.Ward {
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

func importTestCases(fileName string) []TestCase {
	var testCases []TestCase

	bytes, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(bytes, &testCases); err != nil {
		log.Fatal(err)
	}

	return testCases

}

func logWords(words []string) {
	log.Println("Words Count:", len(words))
	for i, word := range words {
		log.Println(i+1, word)
	}
}

func importDictionary(fileName string) *trie.Trie {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	trieDic := trie.NewTrie()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		trieDic.AddWord(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return trieDic
}
