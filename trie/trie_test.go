package trie

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrie_ExtractWord(t *testing.T) {
	trie := NewTrie()

	trie.AddWord("nguyen")
	trie.AddWord("nguyen tri phuong")
	trie.AddWord("tp ho chi minh")

	type args struct {
		name     string
		sentence string
		offset   int
	}

	cases := []struct {
		args args
		want string
	}{
		{
			args: args{
				name:     "case 1",
				sentence: "nguyen",
				offset:   0,
			},
			want: "nguyen",
		},
		{
			args: args{
				name:     "case 2",
				sentence: "nguyen tri phuong",
				offset:   0,
			},
			want: "nguyen",
		},
		{
			args: args{
				sentence: "nguyen tri phuong",
				offset:   6,
			},
			want: "nguyen tri phuong",
		},
		{
			args: args{
				sentence: "nguyen tri phuong, phuong 10, quan 10, tp ho chi minh",
				offset:   6,
			},
			want: "nguyen tri phuong",
		},
		{
			args: args{
				sentence: "tp ho chi minh",
				offset:   0,
			},
			want: "tp ho chi minh",
		},
	}

	for _, c := range cases {
		t.Run(c.args.name, func(t *testing.T) {
			assert.Equal(t, c.want, trie.ExtractWord(c.args.sentence, c.args.offset))
		})
	}
}
