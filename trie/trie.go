package trie

import (
	"address_classification/entity"
	"address_classification/pkg/stringutil"
	"fmt"
	"math"
	"strings"
)

type Trie struct {
	Root *Node
}

type NodeType int

const (
	NodeTypeOther    NodeType = 0
	NodeTypeWard     NodeType = 1
	NodeTypeDistrict NodeType = 2
	NodeTypeProvince NodeType = 3
)

func (nt NodeType) ToString() string {
	switch nt {
	case NodeTypeWard:
		return "Ward"
	case NodeTypeDistrict:
		return "District"
	case NodeTypeProvince:
		return "Province"
	default:
		return "Other"
	}
}

type Node struct {
	Weight   int
	Height   int
	Value    string
	IsEnd    bool
	Type     NodeType
	IDs      []string
	Children map[rune]*Node
}

var skipMap = make(map[string]int)

func NewTrie() *Trie {
	return &Trie{Root: &Node{Children: make(map[rune]*Node)}}
}

func (trie *Trie) setCacheSkip(sentence string, skip int) {
	skipMap[sentence] = skip
}

func (trie *Trie) getCacheSkip(sentence string) int {
	return skipMap[sentence]
}

func (trie *Trie) BuildTrie(words []string) {
	for _, word := range words {
		trie.AddWord(word)
	}
}

func (trie *Trie) AddWord(word string) {
	node := trie.Root

	for _, char := range word {
		//if char == ' ' {
		//	node.IsEnd = true
		//}

		if _, ok := node.Children[char]; !ok {
			node.Children[char] = &Node{Children: make(map[rune]*Node)}
		}
		height := node.Height
		node = node.Children[char]
		node.Value = string(char)
		node.Height = height + 1
	}
	node.IsEnd = true
}

func (trie *Trie) AddWordWithTypeAndID(word string, nodeType NodeType, id string) {
	node := trie.Root

	for _, char := range word {
		if _, ok := node.Children[char]; !ok {
			node.Children[char] = &Node{Children: make(map[rune]*Node)}
		}
		height := node.Height
		node = node.Children[char]
		node.Value = string(char)
		node.Height = height + 1
	}

	node.Type = nodeType
	node.IDs = append(node.IDs, id)

	node.IsEnd = true
}

var wardMap = make(map[string]entity.Ward)
var districtMap = make(map[string]entity.District)
var provinceMap = make(map[string]entity.Province)

func (trie *Trie) BuildTrieWithWards(wards []entity.Ward) {
	for _, ward := range wards {
		wardMap[ward.Code] = ward
		districtMap[ward.DistrictCode] = entity.District{Name: ward.District, Code: ward.DistrictCode}
		provinceMap[ward.ProvinceCode] = entity.Province{Name: ward.Province, Code: ward.ProvinceCode}

		wardName := strings.ToLower(stringutil.RemoveVietnameseAccents(ward.Name))
		trie.AddWordWithTypeAndID(wardName, NodeTypeWard, ward.Code)
	}

	for _, district := range districtMap {
		districtName := strings.ToLower(stringutil.RemoveVietnameseAccents(district.Name))
		trie.AddWordWithTypeAndID(districtName, NodeTypeDistrict, district.Code)
	}

	for _, province := range provinceMap {
		provinceName := strings.ToLower(stringutil.RemoveVietnameseAccents(province.Name))
		trie.AddWordWithTypeAndID(provinceName, NodeTypeProvince, province.Code)
		provinceName = strings.TrimPrefix(provinceName, "thanh pho ")
		provinceName = strings.TrimPrefix(provinceName, "tinh ")
		trie.AddWordWithTypeAndID(provinceName, NodeTypeProvince, province.Code)
	}
}

func (trie *Trie) Print() {
	var dfs func(node *Node, prefix string)
	dfs = func(node *Node, prefix string) {
		if node.IsEnd {
			fmt.Printf("%s, %s %v\n", prefix, node.Type.ToString(), node.IDs) // Print the word when you reach the end of it
		}
		for char, child := range node.Children {
			dfs(child, prefix+string(char)) // Recursively print child nodes
		}
	}
	fmt.Println("------------Start print trie ----------")
	dfs(trie.Root, "")
	fmt.Println("------------End print trie ----------")
}

func (trie *Trie) IsEnd(word string) bool {
	node := trie.Root
	for _, char := range word {
		child, ok := node.Children[char]
		if !ok {
			return false
		}
		node = child
	}
	return node.IsEnd
}

func (trie *Trie) ExtractWord(sentence string, offset int) string {
	node := trie.Root
	for i, char := range sentence {
		if i > offset && node.IsEnd {
			break
		}
		child, ok := node.Children[char]
		if !ok {
			return ""
		}
		node = child
	}
	if !node.IsEnd {
		return ""
	}

	return sentence[:node.Height]
}

func (trie *Trie) ExtractWordWithSkipping(sentence string, offset int) (string, int) {
	var result string
	skip := trie.getCacheSkip(sentence)

	for {
		result = trie.ExtractWord(sentence[skip:], offset)
		if result != "" {
			break
		}

		skip += 1
		if skip >= len(sentence) {
			break
		}
	}

	if skip > 0 {
		trie.setCacheSkip(sentence, skip)
	}

	return result, skip
}

// FindWordsWithPrefix tìm tất cả các từ bắt đầu bằng prefix
func (trie *Trie) FindWordsWithPrefix(prefix string) []string {
	node := trie.searchPrefix(prefix)
	words := []string{}
	if node != nil {
		trie.dfs(node, prefix, &words)
	}
	return words
}

// searchPrefix tìm node chứa tiền tố
func (trie *Trie) searchPrefix(prefix string) *Node {
	node := trie.Root
	for _, char := range prefix {
		if _, ok := node.Children[char]; !ok {
			return nil
		}
		node = node.Children[char]
	}
	return node
}

// dfs thực hiện tìm kiếm theo chiều sâu để lấy các từ hoàn chỉnh
func (trie *Trie) dfs(node *Node, prefix string, words *[]string) {
	if node.IsEnd {
		*words = append(*words, prefix)
	}
	for char, child := range node.Children {
		trie.dfs(child, prefix+string(char), words)
	}
}

// LevenshteinDistance tính khoảng cách Levenshtein giữa hai chuỗi
func LevenshteinDistance(word1, word2 string) int {
	len1 := len(word1)
	len2 := len(word2)
	dp := make([][]int, len1+1)
	for i := range dp {
		dp[i] = make([]int, len2+1)
	}

	for i := 0; i <= len1; i++ {
		for j := 0; j <= len2; j++ {
			if i == 0 {
				dp[i][j] = j
			} else if j == 0 {
				dp[i][j] = i
			} else if word1[i-1] == word2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + int(math.Min(float64(dp[i-1][j]), math.Min(float64(dp[i][j-1]), float64(dp[i-1][j-1]))))
			}
		}
	}
	return dp[len1][len2]
}

// AutoCorrect thực hiện chức năng tự động sửa lỗi
func AutoCorrect(trie *Trie, inputWord string, maxDistance int) []string {
	suggestions := []string{}
	words := trie.FindWordsWithPrefix("") // Tìm tất cả từ trong Trie
	for _, word := range words {
		distance := LevenshteinDistance(inputWord, word)
		if distance <= maxDistance {
			suggestions = append(suggestions, word)
		}
	}
	return suggestions
}
