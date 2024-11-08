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
	if len(provinceIDs) > 0 {
		for _, id := range provinceIDs {
			filterProvinceLocations = append(filterProvinceLocations, locationMap[id])
		}

		sort.Sort(entity.Locations(filterProvinceLocations))

		provinceID = filterProvinceLocations[0].ID
		result = append(result, locationMap[provinceID])

		// decrease word count
		wordsCountMap[locationMap[provinceID].Name] = wordsCountMap[locationMap[provinceID].Name] - 1
	}

	//filter district
	filterDistrictLocations := []entity.Location{}
	if len(districtIDs) == 1 {
		filterDistrictLocations = append(filterDistrictLocations, locationMap[districtIDs[0]])
	} else if len(districtIDs) > 1 {
		// if we have more than 1 district, filter them
		for _, id := range districtIDs {
			district := DistrictMap[id]
			// filter by province and word count
			if district.ProvinceCode == provinceID && wordsCountMap[locationMap[id].Name] > 0 {
				filterDistrictLocations = append(filterDistrictLocations, locationMap[id])
			}
		}

		sort.Sort(entity.Locations(filterDistrictLocations))

		if len(filterDistrictLocations) == 0 {
			filterDistrictLocations = append(filterDistrictLocations, locationMap[districtIDs[0]])
		}
	}

	//var selectedLocation entity.Location
	//if len(filterDistrictLocations) >= 1 {
	//	selectedLocation = filterDistrictLocations[0]
	//}

	//if len(filterDistrictLocations) > 1 {
	//	// case: district with the same name with province
	//	if len(provinceIDs) > 0 {
	//		if locationMap[provinceIDs[0]].Name == filterDistrictLocations[0].Name && wordsCountMap[filterDistrictLocations[0].Name] <= 1 {
	//			selectedLocation = filterDistrictLocations[1]
	//		}
	//	}
	//}

	if len(filterDistrictLocations) > 0 {
		result = append(result, filterDistrictLocations[0])
		wordsCountMap[filterDistrictLocations[0].Name] = wordsCountMap[filterDistrictLocations[0].Name] - 1
	}

	//filter ward
	filterWardIDs := []string{}
	filterWardLocations := []entity.Location{}
	if len(wardIDs) > 0 {
		// if we have more than 1 ward, filter them
		for _, id := range wardIDs {
			ward := WardMap[id]
			// filter ward by province and word count
			if provinceID == ward.ProvinceCode && wordsCountMap[locationMap[id].Name] > 0 {
				filterWardLocations = append(filterWardLocations, locationMap[id])
				filterWardIDs = append(filterWardIDs, id)
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
		// to be improve
		//for _, l := range filterWardLocations {
		//	ward := WardMap[l.ID]
		//	if locationMap[l.ID].Name == locationMap[ward.ProvinceCode].Name {
		//		// remove ward if it's the same name with province
		//		for i, v := range filterWardLocations {
		//			if v.ID == l.ID {
		//				filterWardLocations = append(filterWardLocations[:i], filterWardLocations[i+1:]...)
		//				break
		//			}
		//		}
		//	}
		//}

		result = append(result, filterWardLocations[0])
		wordsCountMap[filterWardLocations[0].Name] = wordsCountMap[filterWardLocations[0].Name] - 1
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
