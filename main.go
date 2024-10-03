package main

import (
	"address_classification/trie"
	"address_classification/trie/parse"
	"bufio"
	"log"
	"os"
	"strings"
)

type Result struct {
	Street   string
	Ward     string
	District string
	City     string
}

func main() {
	trieDic := importDictionary("./assets/example.txt")

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

	ok, words := parse.DynamicParse(input, trieDic)
	if ok {
		logWords(words)
	}

	return result
}

func normalizeInput(input string) string {
	input = strings.ToLower(input)
	return input
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

	//trieDic.AddWord("Hà Nội")

	return trieDic
}

func logWords(words []string) {
	log.Println("Words Count:", len(words))
	for i, word := range words {
		log.Println(i+1, word)
	}

}
