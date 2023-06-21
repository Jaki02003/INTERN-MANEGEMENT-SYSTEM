package model

type InternStudent struct {
	UserName    string  `json:"username"`
	TotalSolved int     `json:"totalsolved"`
	CGPA        float64 `json:"cgpa"`
}
