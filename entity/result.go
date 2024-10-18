package entity

type Result struct {
	Ward     string `json:"ward"`
	District string `json:"district"`
	Province string `json:"province"`
}

func (r *Result) IsComplete() bool {
	return r.Ward != "" && r.District != "" && r.Province != ""
}
