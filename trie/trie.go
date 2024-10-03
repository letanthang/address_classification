package trie

import "fmt"

type Trie struct {
	Root *Node
}

func NewTrie() *Trie {
	return &Trie{Root: &Node{Children: make(map[rune]*Node)}}
}

type Node struct {
	Weight   int
	Height   int
	Value    string
	IsEnd    bool
	Children map[rune]*Node
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

func (trie *Trie) Print() {
	var dfs func(node *Node, prefix string)
	dfs = func(node *Node, prefix string) {
		if node.IsEnd {
			fmt.Println(prefix) // Print the word when you reach the end of it
		}
		for char, child := range node.Children {
			dfs(child, prefix+string(char)) // Recursively print child nodes
		}
	}
	dfs(trie.Root, "")
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
