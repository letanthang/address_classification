package trie

import (
	"address_classification/entity"
	"address_classification/pkg/stringutil"
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Trie struct {
	Root     *Node
	reversed bool
}

type Node struct {
	Weight    int
	Height    int
	Value     string
	IsEnd     bool
	Locations []entity.Location
	Children  map[rune]*Node
}

type WordDistance struct {
	Word     string
	Distance int
}

const (
	HighWeight   = 5
	MediumWeight = 4
	LowWeight    = 3
	LowestWeight = 2
)

var (
	WardMap     = make(map[string]entity.Ward)
	DistrictMap = make(map[string]entity.District)
	ProvinceMap = make(map[string]entity.Province)

	skipMap = make(map[string]int)
)

func NewTrie(reversed bool) *Trie {
	return &Trie{Root: &Node{Children: make(map[rune]*Node), Weight: 0}, reversed: reversed}
}

func (trie *Trie) AddWordWithTypeAndID(word string, locationType entity.LocationType, id string, weight int) {
	if trie.reversed {
		word = stringutil.Reverse(word)
	}

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

func (trie *Trie) BuildTrieWithWards(wards []entity.Ward) {
	name := ""
	for _, ward := range wards {
		// remove prefix for ward, district, and province
		noPrefixWardName := stringutil.RemoveWardPrefix(ward.Name)
		noPrefixDistrictName := stringutil.RemoveDistrictPrefix(ward.District)
		noPrefixProvinceName := stringutil.RemoveProvincePrefix(ward.Province)

		noPrefixNoAccentWardName := stringutil.StandardizeLocation(noPrefixWardName)
		noPrefixNoAccentDistrictName := stringutil.StandardizeLocation(noPrefixDistrictName)

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
				if noPrefixWardName == noPrefixDistrictName || noPrefixNoAccentWardName == noPrefixNoAccentDistrictName {
					trie.AddWordWithTypeAndID(name, entity.LocationTypeWard, ward.Code, LowestWeight)
				} else {
					trie.AddWordWithTypeAndID(name, entity.LocationTypeWard, ward.Code, LowWeight)
				}
			}
			aliases := []string{"p ", "p.", "p. "}
			for _, alias := range aliases {
				trie.AddWordWithTypeAndID(alias+name, entity.LocationTypeWard, ward.Code, MediumWeight)

				if numAlias, ok := stringutil.NumberWardAliasMap[name]; ok {
					trie.AddWordWithTypeAndID(alias+numAlias, entity.LocationTypeWard, ward.Code, MediumWeight)
				}
			}

			alias := "p"
			// case: number ward
			if _, err := strconv.Atoi(name); err == nil {
				trie.AddWordWithTypeAndID(alias+name, entity.LocationTypeWard, ward.Code, MediumWeight)
				if numAlias, ok := stringutil.NumberWardAliasMap[name]; ok {
					trie.AddWordWithTypeAndID(alias+numAlias, entity.LocationTypeWard, ward.Code, MediumWeight)
				}
			}
		}

		if strings.HasPrefix(wardName, "thi tran ") {
			name = strings.TrimPrefix(wardName, "thi tran ")
			if len(name) > 4 {
				if noPrefixWardName == noPrefixDistrictName {
					trie.AddWordWithTypeAndID(name, entity.LocationTypeWard, ward.Code, LowestWeight)
				} else {
					trie.AddWordWithTypeAndID(name, entity.LocationTypeWard, ward.Code, LowWeight)
				}

			}
			alias := []string{"tt", "tt ", "tt.", "tt. ", "t.t ", "t.t. "}
			for _, a := range alias {
				trie.AddWordWithTypeAndID(a+name, entity.LocationTypeWard, ward.Code, MediumWeight)
			}
		}

	}

	for _, district := range DistrictMap {
		noPrefixDistrictName := stringutil.RemoveDistrictPrefix(district.Name)
		noPrefixProvinceName := stringutil.RemoveProvincePrefix(ProvinceMap[district.ProvinceCode].Name)

		districtName := strings.ToLower(stringutil.RemoveVietnameseAccents(district.Name))
		trie.AddWordWithTypeAndID(districtName, entity.LocationTypeDistrict, district.Code, HighWeight)

		if strings.HasPrefix(districtName, "thi xa ") {
			name = strings.TrimPrefix(districtName, "thi xa ")
			if noPrefixDistrictName == noPrefixProvinceName {
				trie.AddWordWithTypeAndID(name, entity.LocationTypeDistrict, district.Code, LowestWeight)
			} else {
				trie.AddWordWithTypeAndID(name, entity.LocationTypeDistrict, district.Code, LowWeight)
			}

			alias := []string{"tx", "tx ", "tx. ", "t.x ", "t.x. "}
			for _, a := range alias {
				trie.AddWordWithTypeAndID(a+name, entity.LocationTypeDistrict, district.Code, MediumWeight)
			}
		}

		// thanh pho my tho
		if strings.HasPrefix(districtName, "thanh pho ") {
			name = strings.TrimPrefix(districtName, "thanh pho ")
			if name != "vinh" {
				trie.AddWordWithTypeAndID(name, entity.LocationTypeDistrict, district.Code, LowWeight)
			}

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

func (trie *Trie) Skip(sentence string) int {
	var (
		result string
		skip   int
	)

	for {
		result, _ = trie.ExtractWord(sentence[skip:], 0)
		if result != "" {
			break
		}

		skip += 1
		if skip >= len(sentence) {
			break
		}
	}

	return skip
}

func (trie *Trie) ExtractWordWithAutoCorrect(word string) (string, WordDistance, *Node) {
	if trie.reversed {
		word = stringutil.Reverse(word)
	}
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
		word = stringutil.RemoveDelimeter(word)
		distances := []WordDistance{}
		words := trie.FindWordsWithPrefix(prefix)

		if len(words) == 0 {
			return "", WordDistance{}, nil
		}

		for _, w := range words {
			distance := LevenshteinDistance(word, w)
			distances = append(distances, WordDistance{Word: w, Distance: distance})
		}

		slices.SortFunc(distances, func(i, j WordDistance) int {
			return cmp.Compare(i.Distance, j.Distance)
		})

		minDT := distances[0].Distance
		minDistances := []WordDistance{}
		for _, dt := range distances {
			if dt.Distance == minDT {
				minDistances = append(minDistances, dt)
			}
		}

		_, targetNode := trie.ExtractWord(minDistances[0].Word, 0)

		return distances[0].Word, distances[0], targetNode

	}

	return "", WordDistance{}, nil
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
				trie.AddWordWithTypeAndID(tname, entity.LocationTypeProvince, provinceCode, LowWeight)
			}
		}
	}
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
