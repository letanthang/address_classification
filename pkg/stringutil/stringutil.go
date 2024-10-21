package stringutil

import (
	"strconv"
	"strings"
)

// Function để loại bỏ dấu tiếng Việt
func RemoveVietnameseAccents(input string) string {
	var output strings.Builder
	for _, char := range input {
		if replacement, found := vietnameseTones[char]; found {
			output.WriteRune(replacement)
		} else {
			output.WriteRune(char)
		}
	}
	return output.String()
}

// Map ký tự có dấu sang không dấu
var vietnameseTones = map[rune]rune{
	'á': 'a', 'à': 'a', 'ả': 'a', 'ã': 'a', 'ạ': 'a', 'ă': 'a', 'ắ': 'a', 'ằ': 'a', 'ẳ': 'a', 'ẵ': 'a', 'ặ': 'a',
	'â': 'a', 'ấ': 'a', 'ầ': 'a', 'ẩ': 'a', 'ẫ': 'a', 'ậ': 'a',
	'é': 'e', 'è': 'e', 'ẻ': 'e', 'ẽ': 'e', 'ẹ': 'e', 'ê': 'e', 'ế': 'e', 'ề': 'e', 'ể': 'e', 'ễ': 'e', 'ệ': 'e',
	'í': 'i', 'ì': 'i', 'ỉ': 'i', 'ĩ': 'i', 'ị': 'i',
	'ó': 'o', 'ò': 'o', 'ỏ': 'o', 'õ': 'o', 'ọ': 'o', 'ô': 'o', 'ố': 'o', 'ồ': 'o', 'ổ': 'o', 'ỗ': 'o', 'ộ': 'o',
	'ơ': 'o', 'ớ': 'o', 'ờ': 'o', 'ở': 'o', 'ỡ': 'o', 'ợ': 'o',
	'ú': 'u', 'ù': 'u', 'ủ': 'u', 'ũ': 'u', 'ụ': 'u', 'ư': 'u', 'ứ': 'u', 'ừ': 'u', 'ử': 'u', 'ữ': 'u', 'ự': 'u',
	'ý': 'y', 'ỳ': 'y', 'ỷ': 'y', 'ỹ': 'y', 'ỵ': 'y',
	'Á': 'A', 'À': 'A', 'Ả': 'A', 'Ã': 'A', 'Ạ': 'A', 'Ă': 'A', 'Ắ': 'A', 'Ằ': 'A', 'Ẳ': 'A', 'Ẵ': 'A', 'Ặ': 'A',
	'Â': 'A', 'Ấ': 'A', 'Ầ': 'A', 'Ẩ': 'A', 'Ẫ': 'A', 'Ậ': 'A',
	'É': 'E', 'È': 'E', 'Ẻ': 'E', 'Ẽ': 'E', 'Ẹ': 'E', 'Ê': 'E', 'Ế': 'E', 'Ề': 'E', 'Ể': 'E', 'Ễ': 'E', 'Ệ': 'E',
	'Í': 'I', 'Ì': 'I', 'Ỉ': 'I', 'Ĩ': 'I', 'Ị': 'I',
	'Ó': 'O', 'Ò': 'O', 'Ỏ': 'O', 'Õ': 'O', 'Ọ': 'O', 'Ô': 'O', 'Ố': 'O', 'Ồ': 'O', 'Ổ': 'O', 'Ỗ': 'O', 'Ộ': 'O',
	'Ơ': 'O', 'Ớ': 'O', 'Ờ': 'O', 'Ở': 'O', 'Ỡ': 'O', 'Ợ': 'O',
	'Ú': 'U', 'Ù': 'U', 'Ủ': 'U', 'Ũ': 'U', 'Ụ': 'U', 'Ư': 'U', 'Ứ': 'U', 'Ừ': 'U', 'Ử': 'U', 'Ữ': 'U', 'Ự': 'U',
	'Ý': 'Y', 'Ỳ': 'Y', 'Ỷ': 'Y', 'Ỹ': 'Y', 'Ỵ': 'Y', 'đ': 'd', 'Đ': 'D',
}
var delimiters = []rune{',', '.', '-', '_', '+'}

func RemoveDelimeter(name string) string {
	result := name
	for _, delimiter := range delimiters {
		result = strings.ReplaceAll(result, string(delimiter), "")
	}

	return result
}

func RemoveWardPrefix(name string) string {
	result := name
	prefixes := []string{"Phường ", "Xã ", "Thị trấn "}
	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			result = strings.TrimPrefix(name, prefix)
		}
	}

	return result
}

func RemoveDistrictPrefix(name string) string {
	result := name
	prefixes := []string{"Quận ", "Huyện ", "Thị xã ", "Thành phố "}
	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			result = strings.TrimPrefix(name, prefix)
		}
	}

	return result
}

func RemoveProvincePrefix(name string) string {
	result := name
	prefixes := []string{"Tỉnh ", "Thành phố "}
	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			result = strings.TrimPrefix(name, prefix)
		}
	}

	return result
}

func IsInteger(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func Reverse(s string) string {
	// Chuyển chuỗi thành slice của rune
	runes := []rune(s)

	// Đảo ngược slice của rune
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	// Chuyển lại thành chuỗi
	return string(runes)
}
