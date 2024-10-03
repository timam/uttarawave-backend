package models

type Customer struct {
	Mobile string `json:"mobile"`
	Name   string `json:"name"`
	Flat   string `json:"flat"`
	House  string `json:"house"`
	Road   string `json:"road"`
	Block  string `json:"block"`
	Area   string `json:"area"`
	City   string `json:"city"`
}
