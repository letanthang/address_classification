package main

import (
	"address_classification/entity"
	"address_classification/pkg/triehelper"
	"address_classification/trie"
	"fmt"
	"log"
)

func main() {
	wards := triehelper.ImportWardDB("./assets/wards.csv")
	trieDic := trie.NewTrie()
	trieDic.BuildTrieWithWards(wards)

	//trieDic.PrintWithPrefix("thanh")
	ok := trieDic.IsEnd("chiem hoa")
	fmt.Println(ok)

	input := []string{
		"TT Tân Bình Huyện Yên Sơn, Tuyên Quang",
		//"284DBis Ng Văn Giáo, P3, Mỹ Tho, T.Giang.",
		//"Nà Làng Phú Bình, Chiêm Hoá, Chiêm Hoá, Tuyên Quang",
		//"59/12 Ng-B-Khiêm, Đa Kao Quận 1, TP. Hồ Chí Minh",
		//"46/8F Trung Chánh 2 Trung Chánh, Hóc Môn, TP. Hồ Chí Minh",
		//"nguyen tri phuong, phuong 10, quan 10, tp ho chi minh",
		//"nguyen tri, phuong 10, quan 1, tp ho chi minh",
		//"nguyen tri, phuong 100, quan 11, tp ho chi minh",
		//"nguyen tri phuong 100 quan 111 tp ho chi minh",
		//"quan 111 tp ho chi minh",
		//"tp ho chí minh quận 2", // missing Q2 in db/trie
		//"p Quảng Thọ,T.P Sầm Swn,TY. thanh Hóa",
	}

	//word := trieDic.ExtractWord(input[0], 17)
	//log.Println(word)
	for _, address := range input {
		result := triehelper.ClassifyAddress(address, trieDic)
		logResult(result)
	}
}

func logResult(result entity.Result) {
	log.Printf("Result : Province %s, District %s, Ward %s\n", result.Province, result.District, result.Ward)
}
