package trie

import (
	"address_classification/entity"
	"address_classification/pkg/stringutil"
	"cmp"
	"fmt"
	"math"
	"slices"
	"sort"
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

const (
	HighWeight   = 5
	MediumWeight = 4
	LowWeight    = 3
)

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

func (trie *Trie) AddWordWithTypeAndID(word string, locationType entity.LocationType, id string, weight int) {
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

	location := entity.Location{LocationType: locationType, ID: id, Name: word, Weight: weight}
	node.Locations = append(node.Locations, location)
	node.IsEnd = true
	node.Weight = weight
}

var WardMap = make(map[string]entity.Ward)
var DistrictMap = make(map[string]entity.District)
var ProvinceMap = make(map[string]entity.Province)

func (trie *Trie) BuildTrieWithWards(wards []entity.Ward) {
	name := ""
	for _, ward := range wards {
		// remove prefix for ward, district, and province
		noPrefixWardName := stringutil.RemoveWardPrefix(ward.Name)
		noPrefixDistrictName := stringutil.RemoveDistrictPrefix(ward.District)
		noPrefixProvinceName := stringutil.RemoveProvincePrefix(ward.Province)

		ward.NoPrefixName = noPrefixWardName
		WardMap[ward.Code] = ward
		DistrictMap[ward.DistrictCode] = entity.District{Name: ward.District, NoPrefixName: noPrefixDistrictName, Code: ward.DistrictCode, ProvinceCode: ward.ProvinceCode}
		ProvinceMap[ward.ProvinceCode] = entity.Province{Name: ward.Province, NoPrefixName: noPrefixProvinceName, Code: ward.ProvinceCode}

		wardName := strings.ToLower(stringutil.RemoveVietnameseAccents(ward.Name))
		trie.AddWordWithTypeAndID(wardName, entity.LocationTypeWard, ward.Code, HighWeight)
		if strings.HasPrefix(wardName, "xa ") {
			name = strings.TrimPrefix(wardName, "xa ")
			// exclude thanh because it's to ambiguous
			if name != "thanh" {
				trie.AddWordWithTypeAndID(name, entity.LocationTypeWard, ward.Code, LowWeight)
			}

			alias := []string{"x", "x ", "x.", "x. "}
			for _, a := range alias {
				trie.AddWordWithTypeAndID(a+name, entity.LocationTypeWard, ward.Code, MediumWeight)
			}
		}

		if strings.HasPrefix(wardName, "phuong ") {
			name = strings.TrimPrefix(wardName, "phuong ")
			if !stringutil.IsInteger(name) && len(name) > 3 {
				trie.AddWordWithTypeAndID(name, entity.LocationTypeWard, ward.Code, LowWeight)
			}
			alias := []string{"p", "p ", "p.", "p. "}
			for _, a := range alias {
				trie.AddWordWithTypeAndID(a+name, entity.LocationTypeWard, ward.Code, MediumWeight)
			}
		}

		if strings.HasPrefix(wardName, "thi tran ") {
			name = strings.TrimPrefix(wardName, "thi tran ")
			if len(name) > 4 {
				trie.AddWordWithTypeAndID(name, entity.LocationTypeWard, ward.Code, LowWeight)
			}
			alias := []string{"tt", "tt ", "tt.", "tt. ", "t.t ", "t.t. "}
			for _, a := range alias {
				trie.AddWordWithTypeAndID(a+name, entity.LocationTypeWard, ward.Code, MediumWeight)
			}
		}

	}

	for _, district := range DistrictMap {
		districtName := strings.ToLower(stringutil.RemoveVietnameseAccents(district.Name))
		trie.AddWordWithTypeAndID(districtName, entity.LocationTypeDistrict, district.Code, HighWeight)

		if strings.HasPrefix(districtName, "thi xa ") {
			name = strings.TrimPrefix(districtName, "thi xa ")
			trie.AddWordWithTypeAndID(name, entity.LocationTypeDistrict, district.Code, LowWeight)
			alias := []string{"tx", "tx ", "tx. ", "t.x ", "t.x. "}
			for _, a := range alias {
				trie.AddWordWithTypeAndID(a+name, entity.LocationTypeDistrict, district.Code, MediumWeight)
			}
		}

		// thanh pho my tho
		if strings.HasPrefix(districtName, "thanh pho ") {
			name = strings.TrimPrefix(districtName, "thanh pho ")
			trie.AddWordWithTypeAndID(name, entity.LocationTypeDistrict, district.Code, LowWeight)
			alias := []string{"tp", "tp ", "tp. ", "t ", "t. "}
			for _, a := range alias {
				trie.AddWordWithTypeAndID(a+name, entity.LocationTypeDistrict, district.Code, MediumWeight)
			}
		}

		if strings.HasPrefix(districtName, "quan ") {
			name = strings.TrimPrefix(districtName, "quan ")
			if !stringutil.IsInteger(name) {
				trie.AddWordWithTypeAndID(name, entity.LocationTypeDistrict, district.Code, LowWeight)
			}
			alias := []string{"q", "q ", "q.", "q. "}
			for _, a := range alias {
				trie.AddWordWithTypeAndID(a+name, entity.LocationTypeDistrict, district.Code, MediumWeight)
			}
		}

		if strings.HasPrefix(districtName, "huyen ") {
			name = strings.TrimPrefix(districtName, "huyen ")
			if name != "thanh hoa" {
				trie.AddWordWithTypeAndID(name, entity.LocationTypeDistrict, district.Code, LowWeight)
			}

			alias := []string{"h ", "h.", "h. "}
			for _, a := range alias {
				trie.AddWordWithTypeAndID(a+name, entity.LocationTypeDistrict, district.Code, MediumWeight)
			}
		}
	}

	for _, province := range ProvinceMap {
		provinceName := strings.ToLower(stringutil.RemoveVietnameseAccents(province.Name))
		trie.addProvinceWithPrefixAlias(provinceName, province.Code)
	}
}

func (trie *Trie) getProvinceAlias(provinceName string) []string {
	alias := strings.ReplaceAll(provinceName, " ", "")
	return []string{provinceName, alias}
}

func (trie *Trie) addProvinceWithPrefixAlias(provinceName, provinceCode string) {
	var (
		trimName  string
		trimNames []string
		prefixes  []string
	)

	trie.AddWordWithTypeAndID(provinceName, entity.LocationTypeProvince, provinceCode, HighWeight)
	if strings.HasPrefix(provinceName, "thanh pho ") {
		trimName = strings.TrimPrefix(provinceName, "thanh pho ")
		prefixes = []string{"", "thanh pho ", "tp", "tp ", "tp.", "tp. ", "t.", "t. ", "t.p", "t.p "}

	}

	if strings.HasPrefix(provinceName, "tinh ") {
		trimName = strings.TrimPrefix(provinceName, "tinh ")
		prefixes = []string{"", "tinh ", "t", "t.", "t. "}
	}

	trimNames = trie.getProvinceAlias(trimName)

	for _, tname := range trimNames {
		for _, prefix := range prefixes {
			if prefix != "" {
				trie.AddWordWithTypeAndID(prefix+tname, entity.LocationTypeProvince, provinceCode, MediumWeight)
			} else {
				trie.AddWordWithTypeAndID(tname, entity.LocationTypeProvince, provinceCode, LowWeight, revesed)
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

	if skip > len(sentence) {
		return "", nil, 0
	}

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

func (trie *Trie) ExtractWordWithAutoCorrect(word string) (string, WordDistance, *Node) {
	node := trie.Root

	for _, char := range word {
		child, ok := node.Children[char]
		if !ok {
			break
		}
		node = child
	}

	prefix := word[:node.Height]
	if node.Height > 2 {
		distances := []WordDistance{}
		words := trie.FindWordsWithPrefix(prefix)

		for _, w := range words {
			distance := LevenshteinDistance(word, w)
			distances = append(distances, WordDistance{Word: w, Distance: distance})
		}

		slices.SortFunc(distances, func(i, j WordDistance) int {
			return cmp.Compare(j.Distance, i.Distance)
		})

		_, targetNode := trie.ExtractWord(distances[0].Word, 0)

		return distances[0].Word, distances[0], targetNode

	}

	return "", WordDistance{}, nil
}

type WordDistance struct {
	Word     string
	Distance int
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

func countWords(words []string) map[string]int {
	result := make(map[string]int)
	for _, word := range words {
		result[word] = result[word] + 1
	}

	return result
}
func FilterLocation(locations []entity.Location, words []string) []entity.Location {
	if len(locations) == 0 {
		return nil
	}

	wordsCountMap := countWords(words)

	result := []entity.Location{}
	locationMap, wardIDs, districtIDs, provinceIDs := entity.Locations(locations).Simplify()

	//filter ward
	filterWardIDs := []string{}
	filterWardLocations := []entity.Location{}
	if len(wardIDs) == 1 {
		filterWardLocations = append(filterWardLocations, locationMap[wardIDs[0]])
	} else if len(wardIDs) > 0 { // if we have more than 1 ward, filter them
		for _, id := range wardIDs {
			ward := WardMap[id]
			if slices.Contains(provinceIDs, ward.ProvinceCode) {
				filterWardLocations = append(filterWardLocations, locationMap[id])
				filterWardIDs = append(filterWardIDs, id)
				sort.Sort(entity.Locations(filterWardLocations))
			}
		}

		if len(filterWardIDs) == 0 {
			filterWardLocations = append(filterWardLocations, locationMap[wardIDs[0]])
		}
	}

	if len(filterWardLocations) > 0 {
		// to be improve
		for _, l := range filterWardLocations {
			ward := WardMap[l.ID]
			if locationMap[l.ID].Name == locationMap[ward.ProvinceCode].Name {
				// remove ward if it's the same name with province
				for i, v := range filterWardLocations {
					if v.ID == l.ID {
						filterWardLocations = append(filterWardLocations[:i], filterWardLocations[i+1:]...)
						break
					}
				}
			}
		}

		result = append(result, filterWardLocations[0])
	}

	//filter district
	filterDistrictLocations := []entity.Location{}
	if len(districtIDs) == 1 {
		filterDistrictLocations = append(filterDistrictLocations, locationMap[districtIDs[0]])
	} else if len(districtIDs) > 1 { // if we have more than 1 district, filter them
		for _, id := range districtIDs {
			district := DistrictMap[id]
			if slices.Contains(provinceIDs, district.ProvinceCode) {
				filterDistrictLocations = append(filterDistrictLocations, locationMap[id])
			}
		}

		sort.Sort(entity.Locations(filterDistrictLocations))

		if len(filterDistrictLocations) == 0 {
			filterDistrictLocations = append(filterDistrictLocations, locationMap[districtIDs[0]])
		}
	}

	var selectedLocation entity.Location
	if len(filterDistrictLocations) >= 1 {
		selectedLocation = filterDistrictLocations[0]
	}

	if len(filterDistrictLocations) > 1 {
		// case: district with the same name with province
		if len(provinceIDs) > 0 {
			if locationMap[provinceIDs[0]].Name == filterDistrictLocations[0].Name && wordsCountMap[filterDistrictLocations[0].Name] <= 1 {
				selectedLocation = filterDistrictLocations[1]
			}
		}
	}

	if len(filterDistrictLocations) > 0 {
		result = append(result, selectedLocation)
	}

	//filter province

	if len(provinceIDs) > 0 {
		if len(provinceIDs) > 0 {
			result = append(result, locationMap[provinceIDs[0]])
		}
	}

	return result
}
