package trie

import (
	"address_classification/entity"
	"math"
	"slices"
	"sort"
)

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

	provinceID := ""
	//filter province
	if len(provinceIDs) > 0 {
		if len(provinceIDs) > 0 {
			provinceID = provinceIDs[0]
			result = append(result, locationMap[provinceID])
		}
	}

	//filter ward
	filterWardIDs := []string{}
	filterWardLocations := []entity.Location{}
	if len(wardIDs) == 1 {
		filterWardLocations = append(filterWardLocations, locationMap[wardIDs[0]])
	} else if len(wardIDs) > 0 {
		// if we have more than 1 ward, filter them
		// filter ward by province
		for _, id := range wardIDs {
			ward := WardMap[id]
			if provinceID == ward.ProvinceCode {
				filterWardLocations = append(filterWardLocations, locationMap[id])
				filterWardIDs = append(filterWardIDs, id)
				sort.Sort(entity.Locations(filterWardLocations))
			}
		}

		// filter wards by district
		if len(filterWardIDs) == 0 {
			for _, id := range wardIDs {
				ward := WardMap[id]
				if slices.Contains(districtIDs, ward.DistrictCode) {
					filterWardLocations = append(filterWardLocations, locationMap[id])
					filterWardIDs = append(filterWardIDs, id)
					sort.Sort(entity.Locations(filterWardLocations))
				}
			}
		}

		// if we can't filter ward by district, get the first one
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

	return result
}
