package entity

import (
	"fmt"
)

type Ward struct {
	Name         string
	Code         string
	Province     string
	ProvinceCode string
	District     string
	DistrictCode string
	Type         string
}

type Province struct {
	Name string
	Code string
}

type District struct {
	Name         string
	Code         string
	ProvinceCode string
}

type LocationType int

const (
	LocationTypeOther    LocationType = 0
	LocationTypeWard     LocationType = 1
	LocationTypeDistrict LocationType = 2
	LocationTypeProvince LocationType = 3
)

func (nt LocationType) ToString() string {
	switch nt {
	case LocationTypeWard:
		return "Ward"
	case LocationTypeDistrict:
		return "District"
	case LocationTypeProvince:
		return "Province"
	default:
		return "Other"
	}
}

type Location struct {
	Name         string
	LocationType LocationType
	ID           string
}

func (l Location) ToString() string {
	return fmt.Sprintf("%s-%s-%s", l.Name, l.LocationType.ToString(), l.ID)
}

type Locations []Location

func (ls Locations) ToString() string {
	result := ""
	for _, l := range ls {
		result += l.ToString() + "|"
	}
	result = "Locations: " + result
	return result
}

func (ls Locations) Simplify() (map[string]Location, []string, []string, []string) {
	var locationMap = make(map[string]Location)

	provinceIDs := []string{}
	districtIDs := []string{}
	wardIDs := []string{}

	for _, l := range ls {
		locationMap[l.ID] = l
	}

	for _, v := range locationMap {
		if v.LocationType == LocationTypeProvince {
			provinceIDs = append(provinceIDs, v.ID)
		} else if v.LocationType == LocationTypeWard {
			wardIDs = append(wardIDs, v.ID)
		} else {
			districtIDs = append(districtIDs, v.ID)
		}
	}

	return locationMap, wardIDs, districtIDs, provinceIDs
}
