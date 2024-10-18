package trie

import (
	"address_classification/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrie_ExtractWord(t *testing.T) {
	trie := NewTrie()
	trie.AddWordWithTypeAndID("nguyen", entity.LocationTypeOther, "1")
	trie.AddWordWithTypeAndID("nguyen tri phuong", entity.LocationTypeOther, "2")
	trie.AddWordWithTypeAndID("tp ho chi minh", entity.LocationTypeProvince, "3")
	trie.AddWordWithTypeAndID("phuong 11", entity.LocationTypeProvince, "3")
	trie.AddWordWithTypeAndID("p. quang tho", entity.LocationTypeWard, "4")
	trie.AddWordWithTypeAndID("p quang tho", entity.LocationTypeWard, "4")

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
		{
			args: args{
				sentence: "p. quang tho",
				offset:   0,
			},
			want: "p. quang tho",
		},
		{
			args: args{
				sentence: "p quang tho hehe",
				offset:   0,
			},
			want: "p quang tho",
		},
	}

	for _, c := range cases {
		t.Run(c.args.name, func(t *testing.T) {
			word, _ := trie.ExtractWord(c.args.sentence, c.args.offset)
			assert.Equal(t, c.want, word)
		})
	}
}

func TestTrie_ExtractWord_WithBuildTrie(t *testing.T) {
	trie := NewTrie()
	trie.AddWordWithTypeAndID("nguyen", entity.LocationTypeOther, "1")
	trie.AddWordWithTypeAndID("nguyen tri phuong", entity.LocationTypeOther, "2")
	trie.AddWordWithTypeAndID("tp ho chi minh", entity.LocationTypeProvince, "3")
	trie.AddWordWithTypeAndID("phuong 11", entity.LocationTypeProvince, "3")
	trie.AddWordWithTypeAndID("p. quang tho", entity.LocationTypeWard, "4")

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
				sentence: "p. quang tho",
				offset:   0,
			},
			want: "p. quang tho",
		},
	}

	for _, c := range cases {
		t.Run(c.args.name, func(t *testing.T) {
			word, _ := trie.ExtractWord(c.args.sentence, c.args.offset)
			assert.Equal(t, c.want, word)
		})
	}
}

func TestTrie_BuildTrieWithWards(t *testing.T) {
	trie := NewTrie()
	wards := []entity.Ward{
		{
			Name:         "Phường Văn Miếu",
			Code:         "00181",
			Province:     "Thành phố Hà Nội",
			ProvinceCode: "01",
			District:     "Quận Đống Đa",
			DistrictCode: "006",
		},
		{
			Name:         "Phường Quốc Tử Giám",
			Code:         "00184",
			Province:     "Thành phố Hà Nội",
			ProvinceCode: "01",
			District:     "Quận Đống Đa",
			DistrictCode: "006",
		},
	}
	trie.BuildTrieWithWards(wards)

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
				sentence: "phuong van mieu",
				offset:   0,
			},
			want: "phuong van mieu",
		},
		{
			args: args{
				sentence: "p. van mieu",
				offset:   0,
			},
			want: "p. van mieu",
		},
	}

	for _, c := range cases {
		t.Run(c.args.name, func(t *testing.T) {
			word, _ := trie.ExtractWord(c.args.sentence, c.args.offset)
			assert.Equal(t, c.want, word)
		})
	}
}
