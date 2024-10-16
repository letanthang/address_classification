package trie

import (
	"address_classification/entity"
	"address_classification/pkg/stringutil"
	"fmt"
	"math"
	"slices"
	"strings"
)

type Trie struct {
	Root *Node
}

type Node struct {
	Weight    int
	Height    int
	Value     string
	IsEnd     bool
	Locations []entity.Location
	Children  map[rune]*Node
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

func (trie *Trie) AddWordWithTypeAndID(word string, locationType entity.LocationType, id string) {
	node := trie.Root

	for _, char := range word {
		child, ok := node.Children[char]
		if !ok {
			node.Children[char] = &Node{Children: make(map[rune]*Node)}
			child = node.Children[char]
		}
		height := node.Height
		child.Value = string(char)
		child.Height = height + 1
		node = child
	}

	location := entity.Location{LocationType: locationType, ID: id, Name: word}
	node.Locations = append(node.Locations, location)
	node.IsEnd = true
}

var WardMap = make(map[string]entity.Ward)
var DistrictMap = make(map[string]entity.District)
var ProvinceMap = make(map[string]entity.Province)

func (trie *Trie) BuildTrieWithWards(wards []entity.Ward) {
	name := ""
	for _, ward := range wards {
		WardMap[ward.Code] = ward
		DistrictMap[ward.DistrictCode] = entity.District{Name: ward.District, Code: ward.DistrictCode, ProvinceCode: ward.ProvinceCode}
		ProvinceMap[ward.ProvinceCode] = entity.Province{Name: ward.Province, Code: ward.ProvinceCode}

		wardName := strings.ToLower(stringutil.RemoveVietnameseAccents(ward.Name))
		trie.AddWordWithTypeAndID(wardName, entity.LocationTypeWard, ward.Code)
		if strings.HasPrefix(wardName, "xa ") {
			name = strings.TrimPrefix(wardName, "xa ")
			// exclude thanh because it's to ambiguous
			if name != "thanh" {
				trie.AddWordWithTypeAndID(name, entity.LocationTypeWard, ward.Code)
			}

			alias := []string{"x ", "x.", "x. "}
			for _, a := range alias {
				trie.AddWordWithTypeAndID(a+name, entity.LocationTypeWard, ward.Code)
			}
		}

		if strings.HasPrefix(wardName, "phuong ") {
			name = strings.TrimPrefix(wardName, "phuong ")
			trie.AddWordWithTypeAndID(name, entity.LocationTypeWard, ward.Code)
			alias := []string{"p ", "p.", "p. "}
			for _, a := range alias {
				trie.AddWordWithTypeAndID(a+name, entity.LocationTypeWard, ward.Code)
			}
		}

		if strings.HasPrefix(wardName, "thi tran ") {
			name = strings.TrimPrefix(wardName, "thi tran ")
			trie.AddWordWithTypeAndID(name, entity.LocationTypeWard, ward.Code)
			alias := []string{"tt ", "tt.", "tt. ", "t.t ", "t.t. "}
			for _, a := range alias {
				trie.AddWordWithTypeAndID(a+name, entity.LocationTypeWard, ward.Code)
			}
		}

	}

	for _, district := range DistrictMap {
		districtName := strings.ToLower(stringutil.RemoveVietnameseAccents(district.Name))
		trie.AddWordWithTypeAndID(districtName, entity.LocationTypeDistrict, district.Code)
		if districtName == "huyen thanh hoa" {
			continue
		}

		if strings.HasPrefix(districtName, "thi xa ") {
			name = strings.TrimPrefix(districtName, "thi xa ")
			trie.AddWordWithTypeAndID(name, entity.LocationTypeDistrict, district.Code)
			alias := []string{"tx ", "tx. ", "t.x ", "t.x. "}
			for _, a := range alias {
				trie.AddWordWithTypeAndID(a+name, entity.LocationTypeDistrict, district.Code)
			}
		}

		if strings.HasPrefix(districtName, "quan ") {
			name = strings.TrimPrefix(districtName, "quan ")
			trie.AddWordWithTypeAndID(name, entity.LocationTypeDistrict, district.Code)
			alias := []string{"q", "q ", "q.", "q. "}
			for _, a := range alias {
				trie.AddWordWithTypeAndID(a+name, entity.LocationTypeDistrict, district.Code)
			}
		}

		if strings.HasPrefix(districtName, "huyen ") {
			name = strings.TrimPrefix(districtName, "huyen ")
			trie.AddWordWithTypeAndID(name, entity.LocationTypeDistrict, district.Code)
			alias := []string{"h ", "h.", "h. "}
			for _, a := range alias {
				trie.AddWordWithTypeAndID(a+name, entity.LocationTypeDistrict, district.Code)
			}
		}
	}

	for _, province := range ProvinceMap {
		provinceName := strings.ToLower(stringutil.RemoveVietnameseAccents(province.Name))
		trie.AddWordWithTypeAndID(provinceName, entity.LocationTypeProvince, province.Code)

		if strings.HasPrefix(provinceName, "thanh pho ") {
			name = strings.TrimPrefix(provinceName, "thanh pho ")
			trie.AddWordWithTypeAndID(name, entity.LocationTypeProvince, province.Code)
			alias := []string{"tp", "tp ", "tp.", "tp. ", "t.", "t. ", "t.p", "t.p "}
			for _, a := range alias {
				trie.AddWordWithTypeAndID(a+name, entity.LocationTypeProvince, province.Code)
			}
		}

		if strings.HasPrefix(provinceName, "tinh ") {
			name = strings.TrimPrefix(provinceName, "tinh ")
			trie.AddWordWithTypeAndID(name, entity.LocationTypeProvince, province.Code)
			alias := []string{"t", "t.", "t. "}
			for _, a := range alias {
				trie.AddWordWithTypeAndID(a+name, entity.LocationTypeProvince, province.Code)
			}
		}
	}
}

func (trie *Trie) Print() {
	var dfs func(node *Node, prefix string)
	dfs = func(node *Node, prefix string) {
		if node.IsEnd {

			fmt.Printf("%s, %s\n", prefix, entity.Locations(node.Locations).ToString()) // Print the word when you reach the end of it
		}
		for char, child := range node.Children {
			dfs(child, prefix+string(char)) // Recursively print child nodes
		}
	}
	fmt.Println("------------Start print trie ----------")
	dfs(trie.Root, "")
	fmt.Println("------------End print trie ----------")
}

func (trie *Trie) PrintWithPrefix(prefix string) {
	node := trie.searchPrefix(prefix)

	var dfs func(node *Node, prefix string)
	dfs = func(node *Node, prefix string) {
		if node.IsEnd {
			fmt.Printf("%s, %s \n", prefix, entity.Locations(node.Locations).ToString()) // Print the word when you reach the end of it
		}
		for char, child := range node.Children {
			dfs(child, prefix+string(char)) // Recursively print child nodes
		}
	}
	fmt.Println("------------Start print trie ----------")
	dfs(node, prefix)
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

func (trie *Trie) ExtractWord(sentence string, offset int) (string, *Node) {
	node := trie.Root
	breakFlag := false

	for i, char := range sentence {
		if i > offset && node.IsEnd {
			breakFlag = true
		}

		child, ok := node.Children[char]
		if !ok {
			if breakFlag {
				break
			}

			return "", nil
		} else {
			if breakFlag && !child.IsEnd {
				break
			}
		}
		node = child
	}
	if !node.IsEnd {
		return "", nil
	}

	return sentence[:node.Height], node
}

func (trie *Trie) ExtractWordWithSkipping(sentence string, offset int) (string, *Node, int) {
	var (
		result string
		node   *Node
	)
	skip := trie.getCacheSkip(sentence)

	for {
		result, node = trie.ExtractWord(sentence[skip:], offset)
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

	return result, node, skip
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

func FilterLocation(locations []entity.Location) []entity.Location {
	result := []entity.Location{}
	locationMap, wardIDs, districtIDs, provinceIDs := entity.Locations(locations).Simplify()

	if len(provinceIDs) > 0 {
		for _, id := range wardIDs {
			ward := WardMap[id]
			if slices.Contains(provinceIDs, ward.ProvinceCode) {
				result = append(result, locationMap[id])
			}
		}

		for _, id := range districtIDs {
			district := DistrictMap[id]
			if slices.Contains(provinceIDs, district.ProvinceCode) {
				result = append(result, locationMap[id])
			}
		}

		result = append(result, locationMap[provinceIDs[0]])
	}

	return result
}
