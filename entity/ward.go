package entity

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
	Name string
	Code string
}
