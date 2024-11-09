package trie

import (
	"address_classification/entity"
	"slices"
	"sort"
	"strings"
)

func FilterLocation(locations []entity.Location, words []string, sentence string) []entity.Location {
	if len(locations) == 0 {
		return nil
	}

	wordsCountMap := countWords(words, sentence)

	result := []entity.Location{}
	locationMap, wardIDs, districtIDs, provinceIDs := entity.Locations(locations).Simplify()

	//filter province
	filterProvinceLocations := []entity.Location{}
	provinceID := ""
	if len(provinceIDs) == 1 {
		filterProvinceLocations = append(filterProvinceLocations, locationMap[provinceIDs[0]])
	} else if len(provinceIDs) > 0 {
		for _, id := range provinceIDs {
			filterProvinceLocations = append(filterProvinceLocations, locationMap[id])
		}

		sort.Sort(entity.Locations(filterProvinceLocations))

	}

	// choose only 1 province
	if len(filterProvinceLocations) > 0 {
		provinceID = filterProvinceLocations[0].ID
		result = append(result, locationMap[provinceID])
		wordsCountMap[locationMap[provinceID].Name] = wordsCountMap[locationMap[provinceID].Name] - 1
	}

	//filter district
	filterDistrictLocations := []entity.Location{}
	districtID := ""
	if len(districtIDs) == 1 {
		filterDistrictLocations = append(filterDistrictLocations, locationMap[districtIDs[0]])
	} else if len(districtIDs) > 1 {
		// if we have more than 1 district, filter them
		for _, id := range districtIDs {
			district := DistrictMap[id]
			// filter by province and word count
			if district.ProvinceCode == provinceID && wordsCountMap[locationMap[id].Name] > 0 {
				wordsCountMap[locationMap[id].Name] = wordsCountMap[locationMap[id].Name] - 1
				filterDistrictLocations = append(filterDistrictLocations, locationMap[id])
			}
		}

		sort.Sort(entity.Locations(filterDistrictLocations))

		if len(filterDistrictLocations) == 0 {
			filterDistrictLocations = append(filterDistrictLocations, locationMap[districtIDs[0]])
		}
	}

	// choose only 1 district
	if len(filterDistrictLocations) > 0 {
		districtID = districtIDs[0]
		result = append(result, filterDistrictLocations[0])
		wordsCountMap[filterDistrictLocations[0].Name] = wordsCountMap[filterDistrictLocations[0].Name] - 1
	}

	//filter ward
	filterWardIDs := []string{}
	filterWardLocations := []entity.Location{}
	if len(wardIDs) > 0 {
		// if we have more than 1 ward, filter them
		// filter by both province and district if possible
		if provinceID != "" && districtID != "" {
			for _, id := range wardIDs {
				ward := WardMap[id]
				// filter ward by province and word count
				if provinceID == ward.ProvinceCode && districtID == ward.DistrictCode && wordsCountMap[locationMap[id].Name] > 0 {
					filterWardLocations = append(filterWardLocations, locationMap[id])
					filterWardIDs = append(filterWardIDs, id)
				}
			}
		}

		// filter wards by province
		if len(filterWardIDs) == 0 {
			for _, id := range wardIDs {
				ward := WardMap[id]
				// filter ward by province and word count
				if provinceID == ward.ProvinceCode && wordsCountMap[locationMap[id].Name] > 0 {
					filterWardLocations = append(filterWardLocations, locationMap[id])
					filterWardIDs = append(filterWardIDs, id)
				}
			}
		}

		// filter wards by district
		if len(filterWardIDs) == 0 {
			for _, id := range wardIDs {
				ward := WardMap[id]
				if slices.Contains(districtIDs, ward.DistrictCode) && wordsCountMap[locationMap[id].Name] > 0 {
					filterWardLocations = append(filterWardLocations, locationMap[id])
					filterWardIDs = append(filterWardIDs, id)
				}
			}
		}

		//if we can't filter ward by district, get the first one and still have word count
		if len(filterWardIDs) == 0 && wordsCountMap[locationMap[wardIDs[0]].Name] > 0 && locationMap[wardIDs[0]].Weight > LowWeight {
			filterWardLocations = append(filterWardLocations, locationMap[wardIDs[0]])
		} else {
			sort.Sort(entity.Locations(filterWardLocations))
		}
	}

	if len(filterWardLocations) > 0 {

		if len(filterWardLocations) == 1 {
			result = append(result, filterWardLocations[0])
			wordsCountMap[filterWardLocations[0].Name] = wordsCountMap[filterWardLocations[0].Name] - 1
		} else {
			if filterWardLocations[0].Weight == filterWardLocations[1].Weight {
				result = append(result, filterWardLocations[1])
				wordsCountMap[filterWardLocations[1].Name] = wordsCountMap[filterWardLocations[1].Name] - 1
			} else {
				result = append(result, filterWardLocations[0])
				wordsCountMap[filterWardLocations[0].Name] = wordsCountMap[filterWardLocations[0].Name] - 1
			}
		}

	}

	return result
}

func countWords(words []string, sentence string) map[string]int {
	result := make(map[string]int)
	for _, word := range words {
		count := findOccurrences(sentence, word)
		result[word] = count
	}

	return result
}

func findOccurrences(big, small string) int {
	count := 0
	start := 0

	// Loop to find all occurrences of `a` in `b`
	for {
		index := strings.Index(big[start:], small)
		if index == -1 {
			break
		}

		// Move start index forward to continue searching
		start = start + index + 1
		count++
	}

	return count
}
