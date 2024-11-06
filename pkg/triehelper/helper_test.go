package triehelper

import (
	"address_classification/trie"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrie_ClassifyAddress(t *testing.T) {
	wards := ImportWardDB("./assets/wards.csv")

	trieTree := trie.NewTrie(false)
	trieTree.BuildTrieWithWards(wards)

	reversedTrie := trie.NewTrie(true)
	reversedTrie.BuildTrieWithWards(wards)

	cases := ImportTestCases("./assets/inputs.json")

	for i, c := range cases {
		t.Run(fmt.Sprintf("unit test :%d", i), func(t *testing.T) {
			result := ClassifyAddress(c.Input, trieTree, reversedTrie)
			assert.Equal(t, c.Output.Ward, result.Ward)
			assert.Equal(t, c.Output.District, result.District)
			assert.Equal(t, c.Output.Province, result.Province)
		})
	}
}
